package messages

import "fmt"

type color string

const (
	reset color = "\033[0m"

	Red          color = "\033[31m"
	Green        color = "\033[32m"
	Yellow       color = "\033[33m"
	Blue         color = "\033[34m"
	Purple       color = "\033[35m"
	Cyan         color = "\033[36m"
	White        color = "\033[37m"
	BrightBlack  color = "\033[90m"
	BrightRed    color = "\033[91m"
	BrightGreen  color = "\033[92m"
	BrightYellow color = "\033[93m"
	BrightBlue   color = "\033[94m"
	BrightPurple color = "\033[95m"
	BrightCyan   color = "\033[96m"
	BrightWhite  color = "\033[97m"

	BlackBg        color = "\033[40m"
	RedBg          color = "\033[41m"
	GreenBg        color = "\033[42m"
	YellowBg       color = "\033[43m"
	BlueBg         color = "\033[44m"
	PurpleBg       color = "\033[45m"
	CyanBg         color = "\033[46m"
	WhiteBg        color = "\033[47m"
	BrightBlackBg  color = "\033[100m"
	BrightRedBg    color = "\033[101m"
	BrightGreenBg  color = "\033[102m"
	BrightYellowBg color = "\033[103m"
	BrightBlueBg   color = "\033[104m"
	BrightPurpleBg color = "\033[105m"
	BrightCyanBg   color = "\033[106m"
	BrightWhiteBg  color = "\033[107m"

	Bold      color = "\033[1m"
	Dim       color = "\033[2m"
	Underline color = "\033[4m"
	Blink     color = "\033[5m"
	Reverse   color = "\033[7m"
	Hidden    color = "\033[8m"
)

func ColorizeStr(msg string, c ...color) string {
	colors := ""
	for _, color := range c {
		colors += string(color)
	}
	return fmt.Sprintf("%s%s%s", colors, msg, reset)
}
