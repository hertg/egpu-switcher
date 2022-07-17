package logger

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	red     = color.New(color.FgRed).SprintFunc()
	green   = color.New(color.FgHiGreen).SprintFunc()
	yellow  = color.New(color.FgHiYellow).SprintFunc()
	magenta = color.New(color.FgHiMagenta).SprintFunc()
	blue    = color.New(color.FgHiBlue).SprintFunc()
)

func print(prefix string, format string, args ...any) (int, error) {
	return fmt.Printf(prefix+" "+format+"\n", args...)
}

func Debug(format string, args ...any) (int, error) {
	return print(magenta("[debug]"), format, args...)
}

func Info(format string, args ...any) (int, error) {
	return print(blue("[info]"), format, args...)
}

func Warn(format string, args ...any) (int, error) {
	return print(yellow("[warn]"), format, args...)
}

func Error(format string, args ...any) (int, error) {
	return print(red("[error]"), format, args...)
}
