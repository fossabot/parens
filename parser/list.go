package parser

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/reflection"
)

// MacroFunc will receive un-evaluated list of s-expressions and the
// current scope. In addition, if the macro was accessed through a name
// the name will be passed as well. If the macro was not accessed by name
// (e.g. was result of another list etc.), name will be empty string.
type MacroFunc func(scope Scope, name string, exprs []Expr) (interface{}, error)

// ListExpr represents a list (i.e., a function call) expression.
type ListExpr struct {
	List []Expr
}

// Eval evaluates each s-exp in the list and then evaluates the list itself
// as an s-exp.
func (le ListExpr) Eval(scope Scope) (interface{}, error) {
	if len(le.List) == 0 {
		return le.List, nil
	}

	val, err := le.List[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if macroFn, ok := val.(MacroFunc); ok {
		var name string
		if sym, ok := le.List[0].(SymbolExpr); ok {
			name = sym.Symbol
		}
		return macroFn(scope, name, le.List[1:])
	}

	args := []interface{}{}
	for i := 1; i < len(le.List); i++ {
		arg, err := le.List[i].Eval(scope)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return reflection.Call(val, args...)
}

func buildListExpr(tokens *tokenQueue) (Expr, error) {
	le := ListExpr{}

	for {
		next := tokens.Token(0)
		if next == nil {
			return nil, ErrEOF
		}

		if next.Type == lexer.RPAREN {
			break
		}

		exp, err := buildExpr(tokens)
		if err != nil {
			return nil, err
		}

		if exp != nil {
			le.List = append(le.List, exp)
		}

	}
	tokens.Pop()
	return le, nil
}
