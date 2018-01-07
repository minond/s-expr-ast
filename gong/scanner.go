package gong

import "fmt"

type tokenId string

type token struct {
	id     tokenId
	lexeme []rune
	offset int
}

const (
	numberToken     tokenId = "numtok"
	stringToken     tokenId = "strtok"
	booleanToken    tokenId = "booltok"
	identifierToken tokenId = "idtok"
	parenToken      tokenId = "partok"
	eofToken        tokenId = "eoftok"
	invalidToken    tokenId = "inltok"

	charNil        = rune(0)
	charOpenParen  = rune('(')
	charCloseParen = rune(')')
	charPeriod     = rune('.')
	charDash       = rune('-')
	charGt         = rune('>')
	charFslash     = rune('/')
	charZero       = rune('0')
	charNine       = rune('9')
	charSpace      = rune(' ')
	charTab        = rune('\t')
	charNewline    = rune('\n')
	charReturn     = rune('\r')

	charA = rune('A')
	charb = rune('b')
	charx = rune('x')
	charz = rune('z')
)

func (t token) String() string {
	return fmt.Sprintf("(%s: `%s`)", t.id, string(t.lexeme))
}

func NewToken(id tokenId, lexeme []rune, offset int) token {
	return token{
		id:     id,
		lexeme: lexeme,
		offset: offset,
	}
}

func NewCharToken(id tokenId, lexeme rune, offset int) token {
	return NewToken(id, []rune{lexeme}, offset)
}

/**
 * This is our lisp's grammar:
 *
 * statement	= expression* ;
 *
 * expression	= "'" primary
 *              | "(" primary* ")" ;
 *
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
		} else if isParen(curr) {
			tokens = append(tokens, NewCharToken(parenToken, curr, pos))
		} else if isNumeric(curr) {
			tok := parseNumeric(chars, pos)
			pos += len(tok.lexeme) - 1
			tokens = append(tokens, tok)
		} else if isIdentifier(curr) {
			tok := parseIdentifier(chars, pos)
			pos += len(tok.lexeme) - 1
			tokens = append(tokens, tok)
		} else {
			tokens = append(tokens, NewCharToken(invalidToken, curr, pos))
		}
	}

	return tokens
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
		r == charGt ||
		r == charFslash
}

func isIdentifier(r rune) bool {
	return isAlpha(r) || isOperator(r)
}

func isIdentifierLike(r rune) bool {
	return isIdentifier(r) || isNumeric(r)
}

func next(chars []rune, curr int) rune {
	if len(chars) < curr+1 {
		return charNil
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

func parseIdentifier(chars []rune, pos int) token {
	return token{
		id:     identifierToken,
		lexeme: takeWhile(chars, pos, isIdentifierLike),
		offset: pos,
	}
}

func parseNumeric(chars []rune, pos int) token {
	return token{
		id:     numberToken,
		lexeme: takeWhile(chars, pos, isNumericLike),
		offset: pos,
	}
}
