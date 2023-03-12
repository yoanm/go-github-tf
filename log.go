package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogOutput(level zerolog.Level, disableAnsi bool) {
	//nolint:exhaustruct // Not useful here
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
	return func(originalMsg interface{}) string {
		var msg string

		if level, ok := originalMsg.(string); ok {
			switch level {
			case "trace":
				msg = colorize("Trace", colorMagenta, disableAnsi)
			case "debug":
				msg = colorize("Debug", colorGreen, disableAnsi)
			case "info":
				msg = colorize("Info", colorBlue, disableAnsi)
			case "warn":
				msg = colorize("Warn", colorYellow, disableAnsi)
			case "error":
				msg = colorize("Error", colorRed, disableAnsi)
			case "fatal":
				msg = colorize(colorize("Fatal", colorRed, disableAnsi), colorBold, disableAnsi)
			case "panic":
				msg = colorize(colorize("Panic", colorRed, disableAnsi), colorBold, disableAnsi)
			default:
				msg = colorize("???", colorBold, disableAnsi)
			}
		} else {
			msg = strings.ToUpper(fmt.Sprintf("%s", originalMsg))
		}

		return msg
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
