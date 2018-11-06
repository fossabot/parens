package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spy16/parens"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	scope := makeGlobalScope()
	scope.Bind("exit", cancel)
	scope.Bind("two", 2)
	scope.Bind("square", func(v float64) float64 {
		return v * v
	})
	interpreter := parens.New(scope)

	if len(os.Args) == 2 {
		_, err := interpreter.ExecuteFile(os.Args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		return
	}

	interpreter.DefaultSource = "<REPL>"
	repl := parens.NewREPL(interpreter)
	repl.Banner = "Welcome to Parens REPL!\nType \"(?)\" for help!"
	repl.Start(ctx)
}
