package main

import (
	"os"

	"github.com/cscoding21/csval/gen"
	"github.com/cscoding21/csval/tests"
)

func main() {
	file := os.Getenv("GOFILE")
	if len(file) == 0 {
		println("csval is only meant to be run from within \"go generate\".  exiting...")
		os.Exit(-1)
	}

	st := tests.FooStruct{}
	println(st)

	gen.Generate()
}
