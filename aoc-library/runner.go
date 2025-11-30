package aoclibrary

import (
	"flag"
	"fmt"
)

func RunMain() {
	var selectedDay = -1

	flag.IntVar(&selectedDay, "day", -1, "The day to run")
	flag.Parse()
	if selectedDay == -1 {
		maxDay := -1
		for dayNo := range registry {
			if dayNo > maxDay {
				maxDay = dayNo
			}
		}
		if maxDay == -1 {
			panic("No solutions at all!")
		}
		selectedDay = maxDay
	}
	solution, hasSolution := registry[selectedDay]
	if !hasSolution {
		panic(fmt.Sprintf("No solution for day %v", selectedDay))
	}
	var context = Context{}
	err := solution.Solve(&context)
	if err != nil {
		context.Eprintf("Final error: %v", err)
	}
}
