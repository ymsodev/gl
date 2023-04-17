package gl

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var exprComparer = cmp.Comparer(exprEqual)

func TestParse(t *testing.T) {
	testCases := []struct {
		input    []*token
		expected []expr
	}{
		{
			[]*token{
				{tokLParen, 0, 0, 1, "("},
				{tokSym, 0, 1, 2, "+"},
				{tokNum, 0, 2, 5, "123"},
				{tokSym, 0, 5, 10, "hello"},
				{tokRParen, 0, 10, 11, ")"},
				{tokEof, 0, 11, 11, ""},
			},
			[]expr{
				newList(
					&token{tokLParen, 0, 0, 1, "("},
					&token{tokRParen, 0, 10, 11, ")"},
					[]expr{
						newAtom(&token{tokSym, 0, 1, 2, "+"}),
						newAtom(&token{tokNum, 0, 2, 5, "123"}),
						newAtom(&token{tokSym, 0, 5, 10, "hello"}),
					},
				),
			},
		},
		{
			[]*token{
				{tokLParen, 0, 0, 1, "("},
				{tokLParen, 0, 1, 2, "("},
				{tokRParen, 0, 2, 3, ")"},
				{tokRParen, 0, 3, 4, ")"},
				{tokEof, 0, 4, 4, ""},
			},
			[]expr{
				newList(
					&token{tokLParen, 0, 0, 1, "("},
					&token{tokRParen, 0, 3, 4, ")"},
					[]expr{
						newList(
							&token{tokLParen, 0, 1, 2, "("},
							&token{tokRParen, 0, 2, 3, ")"},
							[]expr{},
						),
					},
				),
			},
		},
	}
	for i, testCase := range testCases {
		input := testCase.input
		expected := testCase.expected
		output, err := parse(input)
		if err != nil {
			t.Errorf("Test case %d failed with unexpected error: %v", i, err)
		}
		if diff := cmp.Diff(output, expected, exprComparer); diff != "" {
			t.Errorf("Test case %d failed (-want +got):\n%s", i, diff)
		}
	}
}

func exprEqual(e1, e2 expr) bool {
	if reflect.TypeOf(e1) != reflect.TypeOf(e2) {
		return false
	}
	switch e1.(type) {
	case *list:
		return listEqual(e1.(*list), e2.(*list))
	case *atom:
		return atomEqual(e1.(*atom), e2.(*atom))
	default:
		return false
	}
}

func listEqual(l1, l2 *list) bool {
	if !reflect.DeepEqual(l1.lp, l2.lp) {
		return false
	}
	if !reflect.DeepEqual(l1.rp, l2.rp) {
		return false
	}
	if len(l1.items) != len(l2.items) {
		return false
	}
	for i, item1 := range l1.items {
		item2 := l2.items[i]
		if !exprEqual(item1, item2) {
			return false
		}
	}
	return true
}

func atomEqual(a1, a2 *atom) bool {
	return reflect.DeepEqual(a1, a2)
}
