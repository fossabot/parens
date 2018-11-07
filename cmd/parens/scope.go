package main

import (
	"github.com/spy16/parens"
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/stdlib"
)

const version = "1.0.0"

const help = `
Welcome to Parens!

Type (exit) or Ctrl+D or Ctrl+C to exit the REPL.

See "cmd/parens/main.go" in the github repository for
more information.

https://github.com/spy16/parens
`

func makeGlobalScope() *parens.Scope {
	scope := parens.NewScope(nil)
	scope.Bind("parens-version", version)

	scope.Bind("?", func() string {
		return help
	})

	scope.Bind("tokenize", func(src string) ([]lexer.Token, error) {
		return lexer.New(src).Tokens()
	})

	stdlib.RegisterBuiltins(scope)
	return scope
}
