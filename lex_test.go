package gl

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tokenComparer = cmp.Comparer(func(t1, t2 *token) bool {
	return reflect.DeepEqual(t1, t2)
})

func TestLex(t *testing.T) {
	testCases := []struct {
		input    string
		expected []*token
	}{
		{
			"()+-",
			[]*token{
				{tokLParen, 0, 0, 1, "("},
				{tokRParen, 0, 1, 2, ")"},
				{tokSym, 0, 2, 3, "+"},
				{tokSym, 0, 3, 4, "-"},
				{tokEof, 0, 4, 4, ""},
			},
		},
		{
			"1234.5678 hello-world",
			[]*token{
				{tokNum, 0, 0, 9, "1234.5678"},
				{tokId, 0, 10, 21, "hello-world"},
				{tokEof, 0, 21, 21, ""},
			},
		},
		{
			"\"abcde\"",
			[]*token{
				{tokStr, 0, 0, 7, "\"abcde\""},
				{tokEof, 0, 7, 7, ""},
			},
		},
		{
			"; comment\n(+ (	* 2 3)\n -10)",
			[]*token{
				{tokLParen, 1, 10, 11, "("},
				{tokSym, 1, 11, 12, "+"},
				{tokLParen, 1, 13, 14, "("},
				{tokSym, 1, 15, 16, "*"},
				{tokNum, 1, 17, 18, "2"},
				{tokNum, 1, 19, 20, "3"},
				{tokRParen, 1, 20, 21, ")"},
				{tokNum, 2, 23, 26, "-10"},
				{tokRParen, 2, 26, 27, ")"},
				{tokEof, 2, 27, 27, ""},
			},
		},
	}
	for i, testCase := range testCases {
		input := testCase.input
		expected := testCase.expected
		output := lex(input)
		if diff := cmp.Diff(output, expected, tokenComparer); diff != "" {
			t.Errorf("Test case %d failed (-want +got):\n%s", i, diff)
		}
	}
}
