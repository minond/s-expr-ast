package gong

import (
	"fmt"
	"strings"
)

type TokenKind string

type Token struct {
	kind   TokenKind
	value  string
	offset int
}

const (
	EOL = "\n"

	EofToken          TokenKind = "eof"
	EolToken          TokenKind = "eol"
	InvalidToken      TokenKind = "invalid"
	IdentifierToken   TokenKind = "identifier"
	StringToken       TokenKind = "string"
	CharacterToken    TokenKind = "character"
	OperatorToken     TokenKind = "operator"
	BracketToken      TokenKind = "bracket"
	DecimalToken      TokenKind = "decimal"
	IntegerToken      TokenKind = "integer"
	HexNumberToken    TokenKind = "hex"
	BinaryNumberToken TokenKind = "binary"
	BooleanToken      TokenKind = "bool"
)

func (tok Token) String() string {
	if tok.kind == EofToken {
		return fmt.Sprintf(`(eof:%d)`, tok.offset)
	} else if tok.kind == EolToken {
		return fmt.Sprintf(`(eol:%d)`, tok.offset)
	} else {
		return fmt.Sprintf(`(%s:%d:%d "%s")`, tok.kind, tok.offset, len(tok.value), tok.value)
	}
}

func token(kind TokenKind, value string, offset int) Token {
	return Token{
		kind:   kind,
		value:  value,
		offset: offset,
	}
}

func Lex(raw string) []Token {
	i := 0
	letters := strings.Split(raw+EOL, "")
	lettersLen := len(letters)

	var tokens []Token

	for ; i < lettersLen; i++ {
		letter := letters[i]

		if isQuote(letter) {
			token, len := parseQuoted(letters, i, letter)
			tokens = append(tokens, token)
			i += len
		} else if isDigit(letter) {
			token, len := parseNumeric(letters, i)
			tokens = append(tokens, token)
			i += len - 1
		} else if isEol(letter) {
			// Ignore the last EOL we added
			if i+1 == lettersLen {
				continue
			}

			tokens = append(tokens, token(EolToken, "", i))
		} else if isSpace(letter) {
			continue
		} else if word, len := lookaheadWord(letters, i); isBoolean(word) {
			tokens = append(tokens, token(BooleanToken, word, i))
			i += len - 1
		} else if word := lookahead(letters, i, 2); isOperator(word) {
			tokens = append(tokens, token(OperatorToken, word, i))
			i += 1
		} else if word := lookahead(letters, i, 1); isOperator(word) {
			tokens = append(tokens, token(OperatorToken, word, i))
		} else if word := lookahead(letters, i, 1); isBracket(word) {
			tokens = append(tokens, token(BracketToken, word, i))
		} else {
			word, len := lookaheadWord(letters, i)
			tokens = append(tokens, token(IdentifierToken, word, i))
			i += len - 1
		}
	}

	return append(tokens, token(EofToken, "", i))
}

func isBinaryDigit(str string) bool {
	return str == "0" || str == "1"
}

func isHexDigit(str string) bool {
	return isDigit(str) ||
		str == "A" ||
		str == "B" ||
		str == "C" ||
		str == "D" ||
		str == "E" ||
		str == "F"
}

func isDigit(str string) bool {
	return str == "0" ||
		str == "1" ||
		str == "2" ||
		str == "3" ||
		str == "4" ||
		str == "5" ||
		str == "6" ||
		str == "7" ||
		str == "8" ||
		str == "9"
}

func isLetter(str string) bool {
	str = strings.ToLower(str)

	return str == "a" ||
		str == "b" ||
		str == "c" ||
		str == "d" ||
		str == "e" ||
		str == "f" ||
		str == "g" ||
		str == "h" ||
		str == "i" ||
		str == "j" ||
		str == "k" ||
		str == "l" ||
		str == "m" ||
		str == "n" ||
		str == "o" ||
		str == "p" ||
		str == "q" ||
		str == "r" ||
		str == "s" ||
		str == "t" ||
		str == "u" ||
		str == "v" ||
		str == "w" ||
		str == "x" ||
		str == "y" ||
		str == "z"
}

func isAphaNumeric(str string) bool {
	return isDigit(str) || isLetter(str)
}

func isEol(str string) bool {
	return str == "\n" || str == "\r"
}

func isSpace(str string) bool {
	return str == " " || str == "\t" || str == "\n" || str == "\r"
}

func isQuote(str string) bool {
	return str == `"` || str == `'`
}

func isStringQuoteEsc(str string) bool {
	return str == `\`
}

func isBoolean(str string) bool {
	return str == "true" || str == "false"
}

func isBracket(str string) bool {
	return str == "[" ||
		str == "]" ||
		str == "(" ||
		str == ")" ||
		str == "{" ||
		str == "}"
}

func isOperator(str string) bool {
	return str == "+" ||
		str == "-" ||
		str == "*" ||
		str == "^" ||
		str == "/" ||
		str == "&" ||
		str == "." ||
		str == "," ||
		str == ";" ||
		str == "<" ||
		str == ">" ||
		str == "|" ||
		str == "@" ||
		str == "!" ||
		str == ":" ||
		str == "=" ||
		str == "\\" ||
		str == "||" ||
		str == "&&" ||
		str == "::" ||
		str == "->" ||
		str == "//"
}

func lookahead(letters []string, start, k int) string {
	buff := ""

	for i := start; i < len(letters) && k != 0; {
		buff = buff + letters[i]
		k--
		i++
	}

	return buff
}

func lookaheadWord(letters []string, start int) (string, int) {
	buff := ""

	for i := start; i < len(letters); i++ {
		curr := letters[i]

		if isSpace(curr) || isBracket(curr) || isQuote(curr) {
			break
		}

		if isOperator(lookahead(letters, i, 1)) || isOperator(lookahead(letters, i, 2)) {
			break
		}

		buff = buff + curr
	}

	return buff, len(buff)
}

// Parses:
//   StringToken
//   CharacterToken
func parseQuoted(letters []string, start int, closingQuote string) (Token, int) {
	kind := StringToken
	buff := ""
	lettersLen := len(letters)

	for i := start; i < lettersLen; i++ {
		curr := letters[i]
		prev := ""

		if i != 0 {
			prev = letters[i-1]
		}

		// Add to the buffer anything that's not an escape of the last EOL.
		if !isStringQuoteEsc(curr) && i+1 != lettersLen {
			buff = buff + curr
		}

		// If we're not at the very beginning, we're parsing a quote - our
		// closingQuote - and the previous char wasn't an escape, we're done.
		if i != start && isQuote(curr) && !isStringQuoteEsc(prev) && curr == closingQuote {
			if closingQuote == "'" {
				switch len(buff[1 : len(buff)-1]) {
				case 0:
					fallthrough
				case 1:
					kind = CharacterToken
					break

				default:
					kind = InvalidToken
				}
			}

			return token(kind, buff, start), i - start
		}
	}

	return token(InvalidToken, buff, start), len(buff)
}

// Parses:
//   IntegerToken
//   DecimalToken
//   HexNumberToken
//   BinaryNumberToken
func parseNumeric(letters []string, start int) (Token, int) {
	buff := ""
	peek := lookahead(letters, start, 2)
	kind := IntegerToken
	isInt := true

	// Is this a non-decimal representation of a number?
	if peek == "0x" || peek == "0X" {
		kind = HexNumberToken
		buff = peek
		start += 2
	} else if peek == "0b" || peek == "0B" {
		kind = BinaryNumberToken
		buff = peek
		start += 2
	}

	for i := start; i < len(letters); i++ {
		curr := letters[i]

		isRealChar := isDigit(curr) || curr == "."
		isRealKind := kind == IntegerToken || kind == DecimalToken || kind == InvalidToken

		if kind == HexNumberToken && isHexDigit(curr) {
			buff = buff + curr
		} else if kind == BinaryNumberToken && isBinaryDigit(curr) {
			buff = buff + curr
		} else if isRealKind && isRealChar {
			if isInt == true && curr == "." {
				kind = DecimalToken
				isInt = false
			} else if !isDigit(curr) {
				kind = InvalidToken
			}

			buff = buff + curr
		} else {
			return token(kind, buff, start), len(buff)
		}
	}

	return token(InvalidToken, buff, start), len(buff)
}
