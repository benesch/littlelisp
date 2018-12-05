package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	type exp struct {
		tok Token
		str string
	}
	lparen := exp{Lparen, "("}
	rparen := exp{Rparen, ")"}
	testCases := []struct {
		name string
		in   string
		out  []exp
	}{
		{"empty", "", nil},
		{"lone number", "42", []exp{{Number, "42"}}},
		{"lone string", `"42"`, []exp{{String, "42"}}},
		{"lone symbol", `sym`, []exp{{Symbol, "sym"}}},
		{"empty list", "()", []exp{lparen, rparen}},
		{"empty string", `""`, []exp{{String, ""}}},
		{"complicated list", `(1 "2" three (4.5))`, []exp{
			lparen,
			{Number, "1"},
			{String, "2"},
			{Symbol, "three"},
			lparen, {Number, "4.5"}, rparen,
			rparen,
		}},
		{"symbol characters", "!&#$", []exp{{Symbol, "!&#$"}}},
		{"no whitespace", "a(b)3((d)e)", []exp{
			{Symbol, "a"},
			lparen, {Symbol, "b"}, rparen,
			{Number, "3"},
			lparen, lparen, {Symbol, "d"}, rparen, {Symbol, "e"}, rparen,
		}},
		{"bad number", "42.42.42", []exp{{Number, "42.42.42"}}},
		{"bad nesting 1", ")(", []exp{rparen, lparen}},
		{"bad nesting 2", "((()", []exp{lparen, lparen, lparen, rparen}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScanner([]byte(tc.in))
			var out []exp
			for {
				tok, str, err := s.Scan()
				if err != nil {
					t.Fatal(err)
				} else if tok == EOF {
					break
				}
				out = append(out, exp{tok, str})
			}
			require.Equal(t, tc.out, out)
		})
	}
}
