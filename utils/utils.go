package utils

const (
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	RED    = "\033[31m"
	ORANGE = "\033[33m"
	PURPLE = "\033[35m"
	WHITE  = "\033[0m"
)

func Pre(a string, col string) string {
	a = col + a + WHITE
	return a
}

var (
	PWat = Pre("[PRICE WATCHER  ]", GREEN)
	TEng = Pre("[TRADE ENGINE   ]", ORANGE)
	TExc = Pre("[TRADE EXECUTOR ]", PURPLE)
)

func padWithColor(color, s string, length ...int) string {
	l := 40
	if len(length) > 0 {
		l = length[0]
	}

	paddedString := " " + color + " " + s + WHITE
	if len(paddedString) > l {
		return paddedString[:l]
	}

	for len(paddedString) < l {
		paddedString += " "
	}
	return paddedString
}

func PadG(s string, length ...int) string {
	return padWithColor(GREEN, s, length...)
}

func PadB(s string, length ...int) string {
	return padWithColor(BLUE, s, length...)
}

func PadR(s string, length ...int) string {
	return padWithColor(RED, s, length...)
}

func PadO(s string, length ...int) string {
	return padWithColor(ORANGE, s, length...)
}

func PadP(s string, length ...int) string {
	return padWithColor(PURPLE, s, length...)
}
