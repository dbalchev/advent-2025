package aoclibrary

import "fmt"

type Solition interface {
	Solve(context *Context) error
}

var registry map[int]Solition = make(map[int]Solition)

func Register(dayNo int, solution Solition) {
	if dayNo < 0 {
		panic(fmt.Sprintf("Invalid day %v", dayNo))
	}
	_, hasDay := registry[dayNo]
	if hasDay {
		panic(fmt.Sprintf("Day %v already registered", dayNo))
	}
	registry[dayNo] = solution
}
