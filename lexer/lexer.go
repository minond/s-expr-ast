package lexer

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
	EOF = "\n"

	InvalidToken      TokenKind = "invalid"
	StringToken       TokenKind = "string"
	RealNumberToken   TokenKind = "real"
	HexNumberToken    TokenKind = "hex"
	BinaryNumberToken TokenKind = "binary"
	BooleanToken      TokenKind = "bool"
)

func (tok Token) String() string {
	return fmt.Sprintf(`(%s "%s")`, tok.kind, tok.value)
}

func Lex(raw string) []Token {
	i := 0
	parts := strings.Split(raw+EOF, "")
	partsEnd := len(parts)

	var tokens []Token

	// lookahead := func(k int) string {
	// 	buffer := ""
	//
	// 	for j := i; j < partsEnd && k != 0; {
	// 		buffer = buffer + parts[j]
	// 		k--
	// 		j++
	// 	}
	//
	// 	return buffer
	// }

	lookaheadWord := func() (string, int) {
		buffer := ""
		len := 0

		for j := i; j < partsEnd; j++ {
			if isSpace(parts[j]) {
				break
			}

			buffer = buffer + parts[j]
			len++
		}

		return buffer, len
	}

	for ; i < partsEnd; i++ {
		letter := parts[i]
		peek := EOF

		if i+1 < partsEnd {
			peek = parts[i+1]
		}

		if isStringQuote(letter) {
			token, len := parseString(parts, i)
			tokens = append(tokens, token)
			i += len
		} else if isDigit(letter) {
			buffer := letter
			kind := RealNumberToken

			// Type of number?
			if peek == "x" {
				kind = HexNumberToken
				i += 1
			} else if peek == "b" {
				kind = BinaryNumberToken
				i += 1
			}

			for i = i + 1; i < partsEnd; i++ {
				next := parts[i]

				if kind == HexNumberToken && isHexDigit(next) {
					buffer = buffer + next
				} else if kind == BinaryNumberToken && isBinaryDigit(next) {
					buffer = buffer + next
				} else if isDigit(next) {
					buffer = buffer + next
				} else {
					tokens = append(tokens, Token{
						kind:  kind,
						value: buffer,
					})

					break
				}
			}
		} else if isSpace(letter) {
			continue
		} else if word, len := lookaheadWord(); isBoolean(word) {
			i += len
			tokens = append(tokens, Token{
				kind:  BooleanToken,
				value: word,
			})
		} else {
			word, len := lookaheadWord()
			i += len

			tokens = append(tokens, Token{
				kind:  InvalidToken,
				value: word,
			})
		}
	}

	return tokens
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

func parseString(letters []string, start int) (Token, int) {
	buffer := ""

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
				value: buffer,
			}, i - start
		} else if !isStringQuoteEsc(curr) {
			buffer = buffer + curr
		}
	}

	return Token{}, 0
}
