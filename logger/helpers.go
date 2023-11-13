package logger

import (
	"time"
)

const (
	DBLUE   = "\u001b[38;2;0;51;204m"
	DGREEN  = "\033[38;2;0;153;0m"
	DPURPLE = "\033[38;2;102;0;153m"
	GREEN   = "\033[32m"
	BLUE    = "\033[34m"
	RED     = "\033[31m"
	ORANGE  = "\033[33m"
	PURPLE  = "\033[35m"
	YELLOW  = "\033[93m"
	PINK    = "\033[95m"
	WHITE   = "\033[0m"
	BOLD    = "\u001b[1m"
	NORMAL  = "\u001b[0m"
)

var (
	PWat = "[ PRICE  WATCHER ]"
	TEng = "[  TRADE ENGINE  ]"
	TExc = "[ TRADE EXECUTOR ]"
)

func pad(s string) string {
	l := 40

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
	return time.Now().Format(time.RFC1123)
}
