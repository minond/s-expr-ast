package gong

import (
	"errors"
	"fmt"
)

type tokenId string

type token struct {
	id     tokenId
	lexeme []rune
	offset int
	err    error
}

const (
	numberToken     tokenId = "numtok"
	stringToken     tokenId = "strtok"
	booleanToken    tokenId = "booltok"
	identifierToken tokenId = "idtok"
	quoteToken      tokenId = "quotetok"
	openParenToken  tokenId = "oparentok"
	closeParenToken tokenId = "cparentok"
	eofToken        tokenId = "eoftok"
	invalidToken    tokenId = "inltok"

	charNil        = rune(0)
	charEos        = rune(-1)
	charOpenParen  = rune('(')
	charCloseParen = rune(')')
	charPeriod     = rune('.')
	charDash       = rune('-')
	charQuestion   = rune('?')
	charDblQuote   = rune('"')
	charGt         = rune('>')
	charFslash     = rune('/')
	charZero       = rune('0')
	charNine       = rune('9')
	charSpace      = rune(' ')
	charSngQuote   = rune('\'')
	charBslash     = rune('\\')
	charTab        = rune('\t')
	charNewline    = rune('\n')
	charReturn     = rune('\r')

	charA = rune('A')
	charb = rune('b')
	charx = rune('x')
	charz = rune('z')
)

func (t token) String() string {
	if t.id == eofToken {
		return "(eof)"
	} else if t.err != nil {
		return fmt.Sprintf("(ERROR: %s in (%s: `%s`))", t.err, t.id, string(t.lexeme))
	} else {
		return fmt.Sprintf("(%s: `%s`)", t.id, string(t.lexeme))
	}
}

func NewToken(id tokenId, lexeme []rune, offset int) token {
	return token{
		id:     id,
		lexeme: lexeme,
		offset: offset,
		err:    nil,
	}
}

func NewCharToken(id tokenId, lexeme rune, offset int) token {
	return NewToken(id, []rune{lexeme}, offset)
}

/**
 * primary      = NUMBER
 *              | STRING
 *              | BOOLEAN
 *              | IDENTIFIER
 *              | expression ;
 *
 * NUMBER       = 0 - 9
 *              | NUMBER "." NUMBER
 *              | "0x" [ 0 - 9 | a - f | A - F ]*
 *              | "0b" [ 0 - 1 ]* ;
 *
 * STRING       = "\"" ... "\"" ;
 *
 * BOOLEAN      = #f | #t ;
 *
 * IDENTIFIER   = A - z [ A - z | 0 - 9 | - | > | / ]* ;
 *
 */
func Scan(source string) []token {
	var tokens []token
	var curr rune

	chars := []rune(source)
	total := len(chars)
	pos := 0

	for ; pos < total; pos++ {
		curr = chars[pos]

		if isSpace(curr) {
			continue
		} else if curr == charSngQuote {
			tokens = append(tokens, NewCharToken(quoteToken, curr, pos))
		} else if curr == charOpenParen {
			tokens = append(tokens, NewCharToken(openParenToken, curr, pos))
		} else if curr == charCloseParen {
			tokens = append(tokens, NewCharToken(closeParenToken, curr, pos))
		} else if isNumeric(curr) {
			tok := parseNumeric(chars, pos)
			pos += len(tok.lexeme) - 1
			tokens = append(tokens, tok)
		} else if isIdentifier(curr) {
			tok := parseIdentifier(chars, pos)
			pos += len(tok.lexeme) - 1
			tokens = append(tokens, tok)
		} else if curr == charDblQuote {
			tok := parseString(chars, pos)
			pos += len(tok.lexeme) - 1
			tokens = append(tokens, tok)
		} else {
			tokens = append(tokens, NewCharToken(invalidToken, curr, pos))
		}
	}

	return append(tokens, token{id: eofToken})
}

func isNumeric(r rune) bool {
	return r >= charZero && r <= charNine
}

func isNumericLike(r rune) bool {
	return isNumeric(r) ||
		r == charPeriod ||
		r == charx ||
		r == charb
}

func isAlpha(r rune) bool {
	return r >= charA && r <= charz
}

func isAlphaNumeric(r rune) bool {
	return isNumeric(r) || isAlpha(r)
}

func isSpace(r rune) bool {
	return r == charSpace ||
		r == charTab ||
		r == charNewline ||
		r == charReturn
}

func isParen(r rune) bool {
	return r == charOpenParen || r == charCloseParen
}

func isOperator(r rune) bool {
	return r == charDash ||
		r == charQuestion ||
		r == charGt ||
		r == charFslash
}

func isIdentifier(r rune) bool {
	return isAlpha(r) || isOperator(r)
}

func isIdentifierLike(r rune) bool {
	return isIdentifier(r) || isNumeric(r)
}

func isEof(r rune) bool {
	return r == charEos
}

func next(chars []rune, curr int) rune {
	if curr+1 >= len(chars) {
		return charEos
	} else {
		return chars[curr+1]
	}
}

func takeWhile(chars []rune, pos int, f func(r rune) bool) []rune {
	var buff []rune
	curr := chars[pos]

	for f(curr) {
		buff = append(buff, curr)

		if f(next(chars, pos)) {
			pos += 1
			curr = chars[pos]
		} else {
			break
		}
	}

	return buff
}

func parseString(chars []rune, pos int) token {
	var buff []rune
	var err error

	curr := chars[pos]
	peek := charNil

	for {
		peek = next(chars, pos)

		if curr == charEos {
			err = errors.New("Missing closing parentheses")
			break
		} else if curr == charBslash && peek == charDblQuote {
			// Escaped quote, add quote and skip over escaped char
			buff = append(buff, curr)
			buff = append(buff, peek)
			curr = next(chars, pos+1)
			pos += 2
		} else if curr == charDblQuote && len(buff) != 0 {
			// End of the string
			buff = append(buff, curr)
			break
		} else {
			buff = append(buff, curr)
			pos += 1
			curr = peek
		}
	}

	return token{
		id:     stringToken,
		lexeme: buff,
		offset: pos,
		err:    err,
	}
}

func parseIdentifier(chars []rune, pos int) token {
	return token{
		id:     identifierToken,
		lexeme: takeWhile(chars, pos, isIdentifierLike),
		offset: pos,
		err:    nil,
	}
}

func parseNumeric(chars []rune, pos int) token {
	return token{
		id:     numberToken,
		lexeme: takeWhile(chars, pos, isNumericLike),
		offset: pos,
		err:    nil,
	}
}
