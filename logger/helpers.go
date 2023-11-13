package logger

import (
	"time"
)

const (
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	RED    = "\033[31m"
	ORANGE = "\033[33m"
	PURPLE = "\033[35m"
	YELLOW = "\033[93m"
	PINK   = "\033[95m"
	WHITE  = "\033[0m"
)

func Pre(a string, col string) string {
	a = col + a + WHITE
	return a
}

var (
	PWat = "[ PRICE  WATCHER ]"
	TEng = "[  TRADE ENGINE  ]"
	TExc = "[ TRADE EXECUTOR ]"
)

func pad(s string) string {
	l := 45

	paddedString := " " + s
	if len(paddedString) > l {
		return paddedString[:l]
	}

	for len(paddedString) < l {
		paddedString += " "
	}

	return paddedString
}

func GetTime() string {
	return time.Now().Format("15:04:05.000")
}
