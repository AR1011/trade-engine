package logger

import (
	"time"
)

const (
	ColorDarkBlue   = "\u001b[38;2;0;51;204m"
	ColorDarkGreen  = "\033[38;2;0;153;0m"
	ColorDarkPurple = "\033[38;2;102;0;153m"
	ColorGreen      = "\033[32m"
	ColorBlue       = "\033[34m"
	ColorRed        = "\033[31m"
	ColorOrange     = "\033[33m"
	ColorPurple     = "\033[35m"
	ColorYellow     = "\033[93m"
	ColorPink       = "\033[95m"
	ColorWhite      = "\033[0m"
	FontBold        = "\u001b[1m"
	FontNormal      = "\u001b[0m"
)

var (
	PriceWatcher  = "[ PRICE  WATCHER ]"
	TradeEngine   = "[  TRADE ENGINE  ]"
	TradeExecutor = "[ TRADE EXECUTOR ]"
)

func pad(s string) string {
	// pads the string to 40 characters
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
