package lexer

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLParen
	TokenRParen
	TokenNumLit
	TokenBoolLit
	TokenStringLit
	TokenUnitLit
	TokenSymbol
)

type Token struct {
	Type  TokenType
	Val   string
	Float float64
	Bool  bool
}

func (t Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenNumLit:
		return fmt.Sprintf("NUMLIT(%g)", t.Float)
	case TokenBoolLit:
		return fmt.Sprintf("BOOLLIT(%t)", t.Bool)
	case TokenStringLit:
		return fmt.Sprintf("STRINGLIT(%q)", t.Val)
	case TokenUnitLit:
		return "UNITLIT"
	case TokenSymbol:
		return fmt.Sprintf("SYMBOL(%q)", t.Val)
	default:
		return "UNKNOWN"
	}
}

type Lexer struct {
	input []rune
	pos   int
	err   error
}

func NewLexer(input string) *Lexer { return &Lexer{input: []rune(input)} }

func isDelimiter(c rune) bool {
	return c == '(' || c == ')' || c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (l *Lexer) scanNumber() (Token, error) {
	start := l.pos
	for l.pos < len(l.input) && !isDelimiter(l.input[l.pos]) {
		l.pos++
	}
	text := string(l.input[start:l.pos])
	n, err := ParseNumber(text)
	if err != nil {
		return Token{}, fmt.Errorf("Could not parse %s to a number", text)
	}
	return Token{Type: TokenNumLit, Val: text, Float: n}, nil
}

func (l *Lexer) scanString() (Token, error) {
	l.pos++
	var b strings.Builder
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if c == '"' {
			l.pos++
			return Token{Type: TokenStringLit, Val: b.String()}, nil
		}
		if c != '\\' {
			b.WriteRune(c)
			l.pos++
			continue
		}
		l.pos++
		if l.pos >= len(l.input) {
			b.WriteByte('\\')
			break
		}
		switch c = l.input[l.pos]; c {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case 'b':
			b.WriteByte('\b')
		case '"':
			b.WriteByte('"')
		default:
			b.WriteByte('\\')
			b.WriteRune(c)
		}
		l.pos++
	}
	return Token{}, fmt.Errorf("String literal was never ended (Got \"%s)", b.String())
}

func (l *Lexer) scanSymbol() Token {
	start := l.pos
	for l.pos < len(l.input) && !isDelimiter(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenSymbol, Val: string(l.input[start:l.pos])}
}

func (l *Lexer) keywordAt(word string) bool {
	w := []rune(word)
	end := l.pos + len(w)
	if end > len(l.input) || string(l.input[l.pos:end]) != word {
		return false
	}
	return end == len(l.input) || isDelimiter(l.input[end])
}

func (l *Lexer) nextToken() (Token, error) {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
		case '#':
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		default:
			goto ready
		}
	}
	return Token{Type: TokenEOF}, nil
ready:
	c := l.input[l.pos]
	switch c {
	case '(':
		l.pos++
		return Token{Type: TokenLParen}, nil
	case ')':
		l.pos++
		return Token{Type: TokenRParen}, nil
	case '"':
		return l.scanString()
	}
	if (c == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9') || (c >= '0' && c <= '9') {
		return l.scanNumber()
	}
	if l.keywordAt("true") {
		l.pos += 4
		return Token{Type: TokenBoolLit, Bool: true}, nil
	}
	if l.keywordAt("false") {
		l.pos += 5
		return Token{Type: TokenBoolLit, Bool: false}, nil
	}
	if l.keywordAt("unit") {
		l.pos += 4
		return Token{Type: TokenUnitLit}, nil
	}
	return l.scanSymbol(), nil
}

func (l *Lexer) NextToken() Token {
	if l.err != nil {
		return Token{Type: TokenEOF}
	}
	tok, err := l.nextToken()
	l.err = err
	if err != nil {
		return Token{Type: TokenEOF}
	}
	return tok
}

func (l *Lexer) Err() error { return l.err }

func Tokenize(input string) ([]Token, error) {
	l := NewLexer(input)
	tokens := make([]Token, 0, len(input)/3+1)
	left, right := 0, 0
	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		switch tok.Type {
		case TokenLParen:
			left++
		case TokenRParen:
			right++
		case TokenEOF:
			if left > right {
				return nil, errors.New("Unbalanced parantheses: one or more unclosed left parantheses")
			}
			if right > left {
				return nil, errors.New("Unbalanced parantheses: one or more superfluous right parantheses")
			}
			return tokens, nil
		}
	}
}

func ParseNumber(text string) (float64, error) {
	s := strings.ReplaceAll(text, "_", "")
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n, nil
	}
	body, sign := s, 1.0
	if strings.HasPrefix(body, "+") || strings.HasPrefix(body, "-") {
		if body[0] == '-' {
			sign = -1
		}
		body = body[1:]
	}
	if !strings.HasPrefix(body, "0x") && !strings.HasPrefix(body, "0X") {
		return 0, errors.New("not a number")
	}
	body = body[2:]
	mantissa, exponent := body, 0
	if i := strings.IndexAny(body, "pP"); i >= 0 {
		mantissa = body[:i]
		var err error
		exponent, err = strconv.Atoi(body[i+1:])
		if err != nil {
			return 0, err
		}
	}
	parts := strings.Split(mantissa, ".")
	if len(parts) > 2 || (len(parts) == 1 && parts[0] == "") || (len(parts) == 2 && parts[0] == "" && parts[1] == "") {
		return 0, errors.New("malformed hex float")
	}
	value := 0.0
	for _, c := range parts[0] {
		d, ok := hexDigit(c)
		if !ok {
			return 0, errors.New("bad hex digit")
		}
		value = value*16 + d
	}
	if len(parts) == 2 {
		scale := 1.0 / 16.0
		for _, c := range parts[1] {
			d, ok := hexDigit(c)
			if !ok {
				return 0, errors.New("bad hex digit")
			}
			value += d * scale
			scale /= 16
		}
	}
	return sign * math.Ldexp(value, exponent), nil
}

func hexDigit(c rune) (float64, bool) {
	switch {
	case c >= '0' && c <= '9':
		return float64(c - '0'), true
	case c >= 'a' && c <= 'f':
		return float64(c-'a') + 10, true
	case c >= 'A' && c <= 'F':
		return float64(c-'A') + 10, true
	default:
		return 0, false
	}
}
