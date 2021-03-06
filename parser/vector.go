package parser

import (
	"fmt"
	"strings"

	"github.com/spy16/parens/lexer"
)

// VectorExpr represents a vector form.
type VectorExpr struct {
	List []Expr
}

// Eval creates a golang slice.
func (ve VectorExpr) Eval(scope Scope) (interface{}, error) {
	lst := []interface{}{}

	for _, expr := range ve.List {
		val, err := expr.Eval(scope)
		if err != nil {
			return nil, err
		}
		lst = append(lst, val)
	}

	return lst, nil
}

func (ve VectorExpr) String() string {
	strs := []string{}
	for _, expr := range ve.List {
		strs = append(strs, fmt.Sprint(expr))
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}

func buildVectorExpr(tokens *tokenQueue) (Expr, error) {
	ve := VectorExpr{}

	for {
		next := tokens.Token(0)
		if next == nil {
			return nil, ErrEOF
		}
		if next.Type == lexer.RVECT {
			break
		}
		exp, err := buildExpr(tokens)
		if err != nil {
			return nil, err
		}

		if exp != nil {
			ve.List = append(ve.List, exp)
		}
	}
	tokens.Pop()
	return ve, nil
}
