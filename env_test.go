package gl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvSimple(t *testing.T) {
	data := map[string]glObj{
		"hello": glStr{"world"},
		"abc":   glNum{123},
		"true":  glBool{true},
		"false": glBool{false},
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
				{map[string]glObj{"abc": glNum{123}}, nil},
				{map[string]glObj{"hello": glStr{"world"}}, nil},
				{map[string]glObj{"+++": glBool{false}}, nil},
			},
			sym:      "hello",
			expected: "world",
		},
		{
			envs: []*env{
				{map[string]glObj{"abc": glNum{123}}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
				{map[string]glObj{}, nil},
			},
			sym:      "abc",
			expected: 123,
		},
		{
			envs: []*env{
				{map[string]glObj{"+": glStr{"a"}}, nil},
				{map[string]glObj{"+": glStr{"b"}}, nil},
				{map[string]glObj{"+": glStr{"c"}}, nil},
				{map[string]glObj{"+": glStr{"d"}}, nil},
				{map[string]glObj{"+": glStr{"e"}}, nil},
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
