package logger

import "fmt"

func printf(prefix string, format string, args ...any) (int, error) {
	return fmt.Printf(prefix+" "+format, args...)
}

func Debugf(format string, args ...any) (int, error) {
	return printf("[debug]", format, args...)
}

func Infof(format string, args ...any) (int, error) {
	return printf("[info]", format, args...)
}

func Warnf(format string, args ...any) (int, error) {
	return printf("[warn]", format, args...)
}

func Errorf(format string, args ...any) (int, error) {
	return printf("[error]", format, args...)
}
