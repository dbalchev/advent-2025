package aoclibrary

import (
	"fmt"

	"github.com/fatih/color"
)

type Context struct {
}

func (context *Context) Eprintf(format string, a ...any) {
	color.Red(format, a...)
}

func (context *Context) Solution(part any, solution any) {
	color.Yellow("Solution %v: ", part)
	color.Green("%v", solution)
	fmt.Println()
}
