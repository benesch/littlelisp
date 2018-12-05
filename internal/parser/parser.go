package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/benesch/littlelisp/internal/scanner"
	"github.com/cockroachdb/apd"
)

type List struct{ Exprs []Expr }
type Symbol struct{ Sym string }
type String struct{ Str string }
type Number struct{ Num *apd.Decimal }

type Expr interface {
	expr()
	fmt.Stringer
}

func (l *List) expr()   {}
func (s *Symbol) expr() {}
func (s *String) expr() {}
func (s *Number) expr() {}

func (l *List) String() string {
	var strs []string
	for _, expr := range l.Exprs {
		strs = append(strs, expr.String())
	}
	return "(" + strings.Join(strs, " ") + ")"
}

func (s *Symbol) String() string {
	return s.Sym
}

func (s *String) String() string {
	return `"` + s.Str + `"`
}

func (n *Number) String() string {
	return n.Num.String()
}

func Parse(src []byte) ([]Expr, error) {
	var exprs []Expr
	s := scanner.NewScanner(src)
	for {
		expr, err := parseNextExpr(s)
		if err != nil {
			return nil, err
		} else if expr == nil {
			break
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

func parseNextExpr(s *scanner.Scanner) (Expr, error) {
	tok, str, err := s.Scan()
	if err != nil {
		return nil, err
	}
	return parseExpr(s, tok, str)
}

func parseList(s *scanner.Scanner) (*List, error) {
	var exprs []Expr
	for {
		tok, str, err := s.Scan()
		if err != nil {
			return nil, err
		}
		if tok == scanner.Rparen {
			break
		}
		expr, err := parseExpr(s, tok, str)
		if err != nil {
			return nil, err
		} else if expr == nil {
			return nil, errors.New("unterminated list")
		}
		exprs = append(exprs, expr)
	}
	return &List{Exprs: exprs}, nil
}

func parseExpr(s *scanner.Scanner, tok scanner.Token, str string) (Expr, error) {
	switch tok {
	case scanner.Symbol:
		return &Symbol{Sym: str}, nil
	case scanner.String:
		return &String{Str: str}, nil
	case scanner.Number:
		num, _, err := apd.NewFromString(str)
		if err != nil {
			return nil, fmt.Errorf("parsing number %q: %s", str, err)
		}
		return &Number{Num: num}, nil
	case scanner.Lparen:
		return parseList(s)
	case scanner.Rparen:
		return nil, errors.New("unexpected rparen")
	case scanner.EOF:
		return nil, nil
	default:
		panic("unreachable")
	}
}
