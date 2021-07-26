package logger

import "fmt"

type cliColor string

const (
	colorGreen  cliColor = "\033[1;32m"
	colorYellow cliColor = "\033[1;33m"
	colorRed    cliColor = "\033[1;31m"
	colorReset  cliColor = "\033[0m"
)

var (
	success   = color(colorGreen)
	clientErr = color(colorYellow)
	serverErr = color(colorRed)
)

func color(c cliColor) func(text interface{}) string {
	return func(text interface{}) string {
		return fmt.Sprintf("%s%v%s", c, text, colorReset)
	}
}

func statusColor(statusCode int) func(text interface{}) string {
	switch {
	case statusCode >= 500:
		return serverErr
	case statusCode >= 400:
		return clientErr
	default:
		return success
	}
}
