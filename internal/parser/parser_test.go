package parser

import (
	"testing"

	"github.com/cockroachdb/apd"
	"github.com/stretchr/testify/require"
)

func num(s string) Expr {
	num, _, err := apd.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return &Number{Num: num}
}

func str(s string) Expr {
	return &String{Str: s}
}

func sym(s string) Expr {
	return &Symbol{Sym: s}
}

func list(exprs ...Expr) Expr {
	return &List{Exprs: exprs}
}

var noErr string

func TestParser(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  []Expr
		err  string
	}{
		{"empty string", "", nil, noErr},
		{"lone number", "42", []Expr{num("42")}, noErr},
		{"lone string", `"42"`, []Expr{str("42")}, noErr},
		{"lone symbol", `sym`, []Expr{sym("sym")}, noErr},
		{"empty list", "()", []Expr{list()}, noErr},
		{"complicated list", `(1 "2" three (4.5))`, []Expr{
			list(num("1"), str("2"), sym("three"), list(num("4.5"))),
		}, noErr},
		{"multiple exprs", `(+ 1 2)()3(print "done!")`, []Expr{
			list(sym("+"), num("1"), num("2")),
			list(),
			num("3"),
			list(sym("print"), str("done!")),
		}, noErr},
		{"bad number", "42.42.42", nil,
			`parsing number "42.42.42": parse mantissa: 4242.42`},
		{"bad nesting 1", ")(", nil, "unexpected rparen"},
		{"bad nesting 2", "((()", nil, "unterminated list"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exprs, err := Parse([]byte(tc.in))
			if tc.err != noErr && (err == nil || tc.err != err.Error()) {
				t.Fatalf("expected %q error, but got %v", tc.err, err)
			} else if tc.err == noErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			} else if err != nil {
				require.Equal(t, tc.out, exprs)
			}
		})
	}
}
