package aoclibrary

import "fmt"

type Solition interface {
	Solve(context *Context) error
}

var registry map[int]Solition = make(map[int]Solition)

func Register(filename string, solution Solition) {
	dayNo := -1
	n, err := fmt.Sscanf(filename, "day%d", &dayNo)
	if n != 1 {
		panic(fmt.Sprintf("Read %d values instead of 1 from %q", n, filename))
	}
	if err != nil {
		panic(err)
	}
	if dayNo < 0 {
		panic(fmt.Sprintf("Invalid day %v", dayNo))
	}
	_, hasDay := registry[dayNo]
	if hasDay {
		panic(fmt.Sprintf("Day %v already registered", dayNo))
	}
	registry[dayNo] = solution
}
