package gl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvSimple(t *testing.T) {
	data := map[string]any{
		"hello": "world",
		"abc":   123,
		"true":  true,
		"false": false,
	}
	env := newEnv(nil)
	for k, v := range data {
		env.set(k, v)
	}
	for k, v := range data {
		res, err := env.get(k)
		if err != nil {
			t.Error(err)
		} else if !cmp.Equal(res, v) {
			t.Errorf("expected %v for %s, got %v", v, k, res)
		}
	}
}

func TestEnvGet(t *testing.T) {
	testCases := []struct {
		envs     []*env
		sym      string
		expected any
	}{
		{
			envs: []*env{
				{map[string]any{"abc": 123}, nil},
				{map[string]any{"hello": "world"}, nil},
				{map[string]any{"+++": false}, nil},
			},
			sym:      "hello",
			expected: "world",
		},
		{
			envs: []*env{
				{map[string]any{"abc": 123}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
				{map[string]any{}, nil},
			},
			sym:      "abc",
			expected: 123,
		},
		{
			envs: []*env{
				{map[string]any{"+": "a"}, nil},
				{map[string]any{"+": "b"}, nil},
				{map[string]any{"+": "c"}, nil},
				{map[string]any{"+": "d"}, nil},
				{map[string]any{"+": "e"}, nil},
			},
			sym:      "+",
			expected: "e",
		},
	}
	for i, testCase := range testCases {
		var outer *env
		for _, env := range testCase.envs {
			env.outer = outer
			outer = env
		}
		val, err := outer.get(testCase.sym)
		if err != nil {
			t.Errorf("Test %d failed: %v", i, err)
		} else if !cmp.Equal(val, testCase.expected) {
			t.Errorf("Test %d failed: expected %v for %s, got %v",
				i, testCase.expected, testCase.sym, val)
		}
	}
}
