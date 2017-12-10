package gong

import (
	"fmt"
	"strings"
)

type TokenKind string

type Token struct {
	kind  TokenKind
	value string
}

const (
	EOL = "\n"

	EofToken          TokenKind = "eof"
	EolToken          TokenKind = "eol"
	InvalidToken      TokenKind = "invalid"
	IdentifierToken   TokenKind = "identifier"
	StringToken       TokenKind = "string"
	OperatorToken     TokenKind = "operator"
	BracketToken      TokenKind = "bracket"
	RealNumberToken   TokenKind = "real"
	HexNumberToken    TokenKind = "hex"
	BinaryNumberToken TokenKind = "binary"
	BooleanToken      TokenKind = "bool"
)

func (tok Token) String() string {
	if tok.kind == EofToken {
		return "(eof)"
	} else if tok.kind == EolToken {
		return "(eol)"
	} else {
		return fmt.Sprintf(`(%s "%s")`, tok.kind, tok.value)
	}
}

func Lex(raw string) []Token {
	letters := strings.Split(raw+EOL, "")
	lettersLen := len(letters)

	var tokens []Token

	for i := 0; i < lettersLen; i++ {
		letter := letters[i]

		if isStringQuote(letter) {
			token, len := parseString(letters, i)
			tokens = append(tokens, token)
			i += len
		} else if isDigit(letter) {
			token, len := parseNumber(letters, i)
			tokens = append(tokens, token)
			i += len - 1
		} else if isEol(letter) {
			// Ignore the last EOL we added
			if i+1 == lettersLen {
				continue
			}

			tokens = append(tokens, Token{
				kind: EolToken,
			})
		} else if isSpace(letter) {
			continue
		} else if word, len := lookaheadWord(letters, i); isBoolean(word) {
			i += len - 1
			tokens = append(tokens, Token{
				kind:  BooleanToken,
				value: word,
			})
		} else if word := lookahead(letters, i, 2); isOperator(word) {
			i += 1
			tokens = append(tokens, Token{
				kind:  OperatorToken,
				value: word,
			})
		} else if word := lookahead(letters, i, 1); isOperator(word) {
			tokens = append(tokens, Token{
				kind:  OperatorToken,
				value: word,
			})
		} else if word := lookahead(letters, i, 1); isBracket(word) {
			tokens = append(tokens, Token{
				kind:  BracketToken,
				value: word,
			})
		} else {
			word, len := lookaheadWord(letters, i)
			i += len - 1

			tokens = append(tokens, Token{
				kind:  IdentifierToken,
				value: word,
			})
		}
	}

	return append(tokens, Token{
		kind: EofToken,
	})
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

func isStringQuote(str string) bool {
	return str == `"`
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
		str == ";" ||
		str == "@" ||
		str == "=" ||
		str == "\\" ||
		str == "::" ||
		str == "->" ||
		str == "//"
}

func empty() (Token, int) {
	return Token{}, 0
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

		if isSpace(curr) || isBracket(curr) {
			break
		}

		if isOperator(lookahead(letters, i, 1)) || isOperator(lookahead(letters, i, 2)) {
			break
		}

		buff = buff + curr
	}

	return buff, len(buff)
}

func parseString(letters []string, start int) (Token, int) {
	buff := ""

	for i := start + 1; i < len(letters); i++ {
		curr := letters[i]
		prev := ""

		if i != 0 {
			prev = letters[i-1]
		}

		// The i = start + 1 up top takes care of skipping the opening quote,
		// and this below takes handles the closing (unescaped) quote.
		if isStringQuote(curr) && !isStringQuoteEsc(prev) {
			return Token{
				kind:  StringToken,
				value: buff,
			}, i - start
		} else if !isStringQuoteEsc(curr) {
			buff = buff + curr
		}
	}

	return empty()
}

func parseNumber(letters []string, start int) (Token, int) {
	buff := ""
	kind := RealNumberToken
	peek := lookahead(letters, start+1, 1)

	isInt := true

	// Type of number?
	if peek == "x" {
		kind = HexNumberToken
		buff = "0x"
		start += 2
	} else if peek == "b" {
		kind = BinaryNumberToken
		buff = "0b"
		start += 2
	}

	for i := start; i < len(letters); i++ {
		curr := letters[i]

		isRealChar := isDigit(curr) || curr == "."
		isRealKind := kind == RealNumberToken || kind == InvalidToken

		if kind == HexNumberToken && isHexDigit(curr) {
			buff = buff + curr
		} else if kind == BinaryNumberToken && isBinaryDigit(curr) {
			buff = buff + curr
		} else if isRealKind && isRealChar {
			if isInt == true && curr == "." {
				isInt = false
			} else if !isDigit(curr) {
				kind = InvalidToken
			}

			buff = buff + curr
		} else {
			return Token{
				kind:  kind,
				value: buff,
			}, len(buff)
		}
	}

	return empty()
}