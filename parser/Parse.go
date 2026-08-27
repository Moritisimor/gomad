package parser

import (
	"errors"
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() *lexer.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.pos]
}

func (p *Parser) advance() *lexer.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	t := &p.tokens[p.pos]
	p.pos++
	return t
}

func (p *Parser) parseExpr() (expr.Expr, error) {
	tok := p.advance()
	if tok == nil {
		return nil, errors.New("unexpected EOF")
	}

	switch tok.Type {
	case lexer.TokenNumLit:
		return expr.NumLit{Val: tok.Float}, nil
	case lexer.TokenBoolLit:
		return expr.BoolLit{Val: tok.Bool}, nil
	case lexer.TokenStringLit:
		return expr.StringLit{Val: tok.Val}, nil
	case lexer.TokenUnitLit:
		return expr.UnitLit{}, nil
	case lexer.TokenSymbol:
		return expr.Symbol{Val: tok.Val}, nil
	case lexer.TokenLParen:
		return p.parseList()
	case lexer.TokenRParen:
		return nil, errors.New("unexpected ')'")
	case lexer.TokenEOF:
		return nil, errors.New("unexpected EOF")
	}
	return nil, fmt.Errorf("unexpected token: %s", tok)
}

func (p *Parser) parseList() (expr.Expr, error) {
	var items []expr.Expr
	for {
		tok := p.peek()
		if tok == nil || tok.Type == lexer.TokenEOF {
			return nil, errors.New("unexpected EOF in list")
		}
		if tok.Type == lexer.TokenRParen {
			p.advance()
			return expr.List{Val: items}, nil
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
}

func ParseProgram(tokens []lexer.Token) ([]expr.Expr, error) {
	p := NewParser(tokens)
	var forms []expr.Expr
	for {
		tok := p.peek()
		if tok == nil || tok.Type == lexer.TokenEOF {
			break
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		forms = append(forms, e)
	}
	return forms, nil
}

func ParseOne(tokens []lexer.Token) (expr.Expr, error) {
	p := NewParser(tokens)
	return p.parseExpr()
}