package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/benesch/littlelisp/internal/parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	exprs, err := parser.Parse(src)
	if err != nil {
		return err
	}
	for _, expr := range exprs {
		fmt.Println(expr)
	}
	return nil
}
