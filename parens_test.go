package parens_test

import (
	"errors"
	"testing"

	"github.com/spy16/parens"
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func add(a, b float64) float64 {
	return a + b
}

func BenchmarkParens_FunctionCall(suite *testing.B) {
	suite.Run("DirectCall", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add(1, 2)
		}
	})

	sexp, err := parser.Parse("<test>", "(add 1 2)")
	if err != nil {
		suite.Fatalf("failed to parse expression: %s", err)
	}

	suite.Run("CallThroughParens", func(b *testing.B) {
		scope := reflection.NewScope(nil)
		scope.Bind("add", add)

		for i := 0; i < b.N; i++ {
			sexp.Eval(scope)
		}
	})
}

func TestExecute_Success(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockSExp(10, nil), nil)

	res, err := par.Execute("10")
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res, 10)
}

func TestExecute_EvalFailure(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(mockSExp(nil, errors.New("failed")), nil)

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func TestExecute_ParseFailure(t *testing.T) {
	scope := reflection.NewScope(nil)
	par := parens.New(scope)
	par.Parse = mockParseFn(nil, errors.New("failed"))

	res, err := par.Execute("(hello)")
	require.Error(t, err)
	assert.Equal(t, errors.New("failed"), err)
	assert.Nil(t, res)
}

func mockSExp(v interface{}, err error) parser.SExp {
	return sexpMock(func(scope *reflection.Scope) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}

func mockParseFn(sexp parser.SExp, err error) parens.ParseFn {
	return func(name, src string) (parser.SExp, error) {
		if err != nil {
			return nil, err
		}
		return sexp, nil
	}
}

type sexpMock func(scope *reflection.Scope) (interface{}, error)

func (sm sexpMock) Eval(scope *reflection.Scope) (interface{}, error) {
	return sm(scope)
}
