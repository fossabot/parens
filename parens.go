package parens

import (
	"github.com/spy16/parens/parser"
	"github.com/spy16/parens/reflection"
)

// New initializes new parens LISP interpreter with given env.
func New(scope *reflection.Scope) *Interpreter {
	return &Interpreter{
		Scope: scope,
		Parse: parser.Parse,
	}
}

// ParseFn is responsible for tokenizing and building SExp out of tokens.
type ParseFn func(src string) (parser.SExp, error)

// Interpreter represents the LISP interpreter instance. You can provide
// your own implementations of ParseFn to extend the interpreter.
type Interpreter struct {
	Scope *reflection.Scope

	// Parse is used to build SExp/AST from source.
	Parse ParseFn
}

// Execute tokenizes, parses and executes the given LISP code.
func (parens *Interpreter) Execute(src string) (interface{}, error) {
	sexp, err := parens.Parse(src)
	if err != nil {
		return nil, err
	}

	res, err := sexp.Eval(parens.Scope)
	if err != nil {
		return nil, err
	}

	return res, nil
}
