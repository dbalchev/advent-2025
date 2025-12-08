package aoclibrary

import (
	"flag"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

var nameToLevel = map[string]slog.Level{
	"INFO":  slog.LevelInfo,
	"DEBUG": slog.LevelDebug,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func parseLevel(strLevel string) slog.Level {

	normalLevel, existing := nameToLevel[strings.ToUpper(strLevel)]
	if existing {
		return normalLevel
	}
	parsedLevel, err := strconv.Atoi(strLevel)
	if err != nil {
		panic(err)
	}
	return slog.Level(parsedLevel)
}

func RunMain() {
	var selectedDay = -1
	var logLevelStr = "INFO"
	var generalArg string

	flag.IntVar(&selectedDay, "day", -1, "The day to run")
	flag.StringVar(&logLevelStr, "log", "INFO", "Log level to use")
	flag.StringVar(&generalArg, "general", "", "a general arg to be used")
	flag.Parse()
	logLevel := parseLevel(logLevelStr)
	slog.SetLogLoggerLevel(logLevel)
	slog.Info("Setting log level", "log level", logLevel)
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
	var context = Context{
		generalArg: generalArg,
	}
	err := solution.Solve(&context)
	if err != nil {
		context.Eprintf("Final error: %v", err)
	}
}
