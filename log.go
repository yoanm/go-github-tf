package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogOutput(level zerolog.Level, disableAnsi bool) {
	outputLogger := zerolog.ConsoleWriter{
		Out: os.Stderr,
		// Remove time header: quite useless for console
		FormatTimestamp: func(i interface{}) string {
			return ""
		},
	}
	outputLogger.FormatLevel = formatLevel(disableAnsi)
	outputLogger.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s", i)
	}
	log.Logger = log.Output(outputLogger)

	zerolog.SetGlobalLevel(level)
}

func formatLevel(disableAnsi bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = colorize("Trace", colorMagenta, disableAnsi)
			case "debug":
				l = colorize("Debug", colorGreen, disableAnsi)
			case "info":
				l = colorize("Info", colorBlue, disableAnsi)
			case "warn":
				l = colorize("Warn", colorYellow, disableAnsi)
			case "error":
				l = colorize("Error", colorRed, disableAnsi)
			case "fatal":
				l = colorize(colorize("Fatal", colorRed, disableAnsi), colorBold, disableAnsi)
			case "panic":
				l = colorize(colorize("Panic", colorRed, disableAnsi), colorBold, disableAnsi)
			default:
				l = colorize("???", colorBold, disableAnsi)
			}
		} else {
			l = strings.ToUpper(fmt.Sprintf("%s", i))
		}
		return l
	}
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)
