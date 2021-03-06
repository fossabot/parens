package stdlib

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/parens"

	"github.com/k0kubun/pp"
	"github.com/spy16/parens/parser"
)

var core = []mapEntry{
	// logical constants
	entry("true", true,
		"Represents logical true",
	),
	entry("false", false,
		"Represents logical false",
	),
	entry("nil", false,
		"Represents logical false. Same as false",
	),

	// core macros
	entry("do", parser.MacroFunc(Do),
		"Usage: (do expr1 expr2 ...)",
	),
	entry("label", parser.MacroFunc(Label),
		"Usage: (label <symbol> expr)",
	),
	entry("global", parser.MacroFunc(Global),
		"Usage: (global <symbol> expr)",
	),
	entry("cond", parser.MacroFunc(Conditional),
		"Usage: (cond (test1 action1) (test2 action2)...)",
	),
	entry("let", parser.MacroFunc(Let),
		"Usage: (let expr1 expr2 ...)",
	),
	entry("inspect", parser.MacroFunc(Inspect),
		"Usage: (inspect expr)",
	),
	entry("lambda", parser.MacroFunc(Lambda),
		"Defines a lambda.",
		"Usage: (lambda (params) body)",
		"where params: a list of symbols",
		"      body  : one or more s-expressions",
	),
	entry("defn", parser.MacroFunc(Defn),
		"Defines a named function",
		"Usage: (defn <name> [params] body)",
	),
	entry("doc", parser.MacroFunc(Doc),
		"Displays documentation for given symbol if available.",
		"Usage: (doc <symbol>)",
	),
	entry("dump-scope", parser.MacroFunc(dumpScope),
		"Formats and displays the entire scope",
	),
	entry("->", parser.MacroFunc(ThreadFirst)),
	entry("->>", parser.MacroFunc(ThreadLast)),

	// core functions
	entry("type", reflect.TypeOf),
}

// ThreadFirst macro appends first evaluation result as first argument of next function
// call.
func ThreadFirst(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	return thread(true, scope, name, exprs)
}

// ThreadLast macro appends first evaluation result as last argument of next function
// call.
func ThreadLast(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	return thread(false, scope, name, exprs)
}

func thread(first bool, scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("at-least 1 argument required")
	}

	var result interface{}
	for i := 1; i < len(exprs); i++ {
		lst, ok := exprs[i].(parser.ListExpr)
		if !ok {
			return nil, fmt.Errorf("argument %d must be a function call, not '%s'", i, reflect.TypeOf(exprs[i]))
		}

		val, err := exprs[i-1].Eval(scope)
		if err != nil {
			return nil, err
		}
		res := anyExpr{val: val}

		nextCall := parser.ListExpr{List: []parser.Expr{lst.List[0]}}

		if first {
			nextCall.List = append(nextCall.List, res)
			nextCall.List = append(nextCall.List, lst.List[1:]...)
		} else {
			nextCall.List = append(nextCall.List, lst.List[1:]...)
			nextCall.List = append(nextCall.List, res)
		}

		result, err = nextCall.Eval(scope)
		if err != nil {
			return nil, err
		}
		exprs[i] = anyExpr{val: result}
	}

	return result, nil
}

type anyExpr struct {
	val interface{}
}

func (ae anyExpr) Eval(scope parser.Scope) (interface{}, error) {
	return ae.val, nil
}

// Doc shows doc string associated with a symbol. If not found, returns a message.
func Doc(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("exactly 1 argument required, got %d", len(exprs))
	}

	sym, ok := exprs[0].(parser.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("argument must be a Symbol, not '%s'", reflect.TypeOf(exprs[0]))
	}

	val, err := sym.Eval(scope)
	if err != nil {
		return nil, err
	}

	docStr := scope.Doc(sym.Symbol)
	if len(strings.TrimSpace(docStr)) == 0 {
		docStr = fmt.Sprintf("No documentation available for '%s'", sym.Symbol)
	}

	docStr = fmt.Sprintf("%s\n\nGo Type: %s", docStr, reflect.TypeOf(val))
	return docStr, nil
}

// Defn macro is for defining named functions. It defines a lambda and binds it with
// the given name into the scope.
func Defn(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) < 3 {
		return nil, fmt.Errorf("3 or more arguments required, got %d", len(exprs))
	}

	sym, ok := exprs[0].(parser.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("first argument must be symbol, not '%s'", reflect.TypeOf(exprs[0]))
	}

	lambda, err := Lambda(scope, name, exprs[1:])
	if err != nil {
		return nil, err
	}

	scope.Bind(sym.Symbol, lambda)
	return sym.Symbol, nil
}

// Lambda macro is for defining lambdas. (lambda (params) body)
func Lambda(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) < 2 {
		return nil, errors.New("at-least two arguments required")
	}

	paramList, ok := exprs[0].(parser.VectorExpr)
	if !ok {
		return nil, fmt.Errorf("first argument must be list of symbols, not '%s'", reflect.TypeOf(exprs[0]))
	}

	params := []string{}
	for _, entry := range paramList.List {
		sym, ok := entry.(parser.SymbolExpr)
		if !ok {
			return nil, fmt.Errorf("param list must contain symbols, not '%s'", reflect.TypeOf(entry))
		}

		params = append(params, sym.Symbol)
	}

	lambdaFunc := func(args ...interface{}) interface{} {
		if len(params) != len(args) {
			panic(fmt.Errorf("requires %d arguments, got %d", len(params), len(args)))
		}

		localScope := parens.NewScope(scope)
		for i := range params {
			localScope.Bind(params[i], args[i])
		}

		val, err := Do(localScope, "", exprs[1:])
		if err != nil {
			panic(err)
		}

		return val
	}

	return lambdaFunc, nil
}

// Do executes all s-exps one by one and returns the result of last evaluation.
func Do(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	var val interface{}
	var err error
	for _, expr := range exprs {
		val, err = expr.Eval(scope)
		if err != nil {
			return nil, err
		}

	}
	return val, nil
}

// Let creates a new sub-scope from the global scope and executes all the
// exprs inside the new scope. Once the Let block ends, all the names bound
// will be removed. In other words, Let is a Do with local scope.
func Let(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	localScope := parens.NewScope(scope)

	return Do(localScope, name, exprs)
}

// Conditional is commonly know LISP (cond (test1 act1)...) construct.
// Tests can be any exressions that evaluate to non-nil and non-false
// value.
func Conditional(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	lists := []parser.ListExpr{}
	for _, exp := range exprs {
		listExp, ok := exp.(parser.ListExpr)

		if !ok {
			return nil, errors.New("all arguments must be lists")
		}
		if len(listExp.List) != 2 {
			return nil, errors.New("each argument must be of the form (test action)")
		}
		lists = append(lists, listExp)
	}

	for _, list := range lists {
		testResult, err := list.List[0].Eval(scope)
		if err != nil {
			return nil, err
		}

		if testResult == nil {
			continue
		}

		if resultBool, ok := testResult.(bool); ok && resultBool == false {
			continue
		}

		return list.List[1].Eval(scope)
	}

	return nil, nil
}

// Label binds the result of evaluating second argument to the symbol passed in as
// first argument in the current scope.
func Label(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	return labelInScope(scope, name, exprs)
}

// Global binds the result of evaluating second argument to the symbol passed in as
// first argument in the global scope.
func Global(scope parser.Scope, name string, exprs []parser.Expr) (interface{}, error) {
	return labelInScope(scope.Root(), name, exprs)
}

// Inspect dumps the exprs in a formatted manner.
func Inspect(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	pp.Println(exprs)
	return nil, nil
}

func dumpScope(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	return fmt.Sprint(scope), nil
}

func labelInScope(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("expecting symbol and a value")
	}
	symbol, ok := exprs[0].(parser.SymbolExpr)
	if !ok {
		return nil, fmt.Errorf("argument 1 must be a symbol, not '%s'", reflect.TypeOf(exprs[0]).String())
	}

	val, err := exprs[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	scope.Bind(symbol.Symbol, val)

	return val, nil
}
