package parser

import (
	"fmt"
)

// QuoteExpr implements the quote-literal form.
type QuoteExpr struct {
	expr Expr
}

// Eval returns the expression itself without evaluating it.
func (qe QuoteExpr) Eval(scope Scope) (interface{}, error) {
	return qe.expr, nil
}

// UnquoteEval unquotes and evaluates the underlying expression.
func (qe QuoteExpr) UnquoteEval(scope Scope) (interface{}, error) {
	return qe.expr.Eval(scope)
}

func (qe QuoteExpr) String() string {
	return fmt.Sprintf("'%s", qe.expr)
}
