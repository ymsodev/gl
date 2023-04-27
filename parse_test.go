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
		expected []glObj
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
			[]glObj{
				glList{
					[]glObj{
						glSym{"+"},
						glNum{123},
						glStr{"hello"},
					},
				},
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
			[]glObj{
				glList{
					[]glObj{
						glList{
							[]glObj{},
						},
					},
				},
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

func exprEqual(obj1, obj2 glObj) bool {
	if reflect.TypeOf(obj1) != reflect.TypeOf(obj2) {
		return false
	}
	switch obj1.(type) {
	case glList:
		return listEqual(obj1.(glList), obj2.(glList))
	default:
		return reflect.DeepEqual(obj1, obj2)
	}
}

func listEqual(l1, l2 glList) bool {
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
