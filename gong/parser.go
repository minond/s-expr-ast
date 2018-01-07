package gong

import (
	"fmt"
	"strings"
)

type expressionKind string

type parser struct {
	tokens []token
	pos    int
}

type expression struct {
	kind  expressionKind
	exprs []expression
	token token
}

const (
	sexpression expressionKind = "sexpression"
	quote       expressionKind = "quote"
	atom        expressionKind = "atom"
)

func (e expression) String() string {
	switch e.kind {
	case atom:
		return fmt.Sprintf("%s", string(e.token.lexeme))

	case quote:
		return fmt.Sprintf("'%s", e.exprs[0])

	case sexpression:
		var content []string

		for _, expr := range e.exprs {
			content = append(content, fmt.Sprintf("%s", expr))
		}

		return fmt.Sprintf("(%s)", strings.Join(content, " "))

	default:
		return "(ERROR)"
	}
}

/**
 * This is our lisp's grammar:
 *
 * MAIN         = expr* ;
 *
 * expression	= "'" primary
 *              | sexpr ;
 *
 * sexpr        = "(" primary* ")" ;
 *
 * primary      = NUMBER
 *              | STRING
 *              | BOOLEAN
 *              | IDENTIFIER
 *              | sexpr ;
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
func Parse(tokens []token) []expression {
	var expressions []expression

	p := parser{
		tokens: tokens,
		pos:    0,
	}

	for !p.done() {
		expressions = append(expressions, p.expression())
	}

	return expressions
}

func (p *parser) expression() expression {
	if p.matches(quoteToken) {
		p.eat()
		return expression{
			kind:  quote,
			exprs: []expression{p.primary()},
		}
	} else {
		return p.sexpr()
	}
}

func (p *parser) sexpr() expression {
	expr := expression{
		kind: sexpression,
	}

	p.expect(openParenToken)

	if p.matches(closeParenToken) {
		p.eat()
		return expr
	}

	for {
		if p.matches(quoteToken) {
			expr.exprs = append(expr.exprs, p.expression())
		} else {
			expr.exprs = append(expr.exprs, p.primary())
		}

		if p.matches(closeParenToken) {
			p.eat()
			break
		}
	}

	return expr
}

func (p *parser) primary() expression {
	switch {
	case p.matches(numberToken):
		fallthrough
	case p.matches(stringToken):
		fallthrough
	case p.matches(booleanToken):
		fallthrough
	case p.matches(identifierToken):
		return expression{
			kind:  atom,
			token: p.eat(),
		}

	case p.matches(openParenToken):
		return p.sexpr()

	default:
		panic("Unkown token")
		return expression{}
	}
}

func (p parser) done() bool {
	return p.pos >= len(p.tokens) || p.tokens[p.pos].id == eofToken
}

func (p parser) peek() token {
	if p.done() {
		return token{id: eofToken}
	} else {
		return p.tokens[p.pos]
	}
}

func (p *parser) eat() token {
	curr := p.peek()
	p.pos += 1
	return curr
}

func (p *parser) expect(id tokenId) error {
	if p.peek().id != id {
		return fmt.Errorf("Expecting %s but found %s", id, p.peek().id)
	} else {
		p.pos += 1
		return nil
	}
}

func (p parser) matches(ids ...tokenId) bool {
	curr := p.peek()

	for _, id := range ids {
		if curr.id == id {
			return true
		}
	}

	return false
}
