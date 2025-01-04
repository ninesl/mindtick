package messages

import "fmt"

type color string

const (
	reset  color = "\033[0m"
	red    color = "\033[31m"
	green  color = "\033[32m"
	yellow color = "\033[33m"
	blue   color = "\033[34m"
	purple color = "\033[35m"
	cyan   color = "\033[36m"
	white  color = "\033[37m"

	blackBg  color = "\033[40m"
	redBg    color = "\033[41m"
	greenBg  color = "\033[42m"
	yellowBg color = "\033[43m"
	blueBg   color = "\033[44m"
	purpleBg color = "\033[45m"
	cyanBg   color = "\033[46m"
	whiteBg  color = "\033[47m"

	bold      color = "\033[1m"
	dim       color = "\033[2m"
	underline color = "\033[4m"
	blink     color = "\033[5m"
	reverse   color = "\033[7m"
	hidden    color = "\033[8m"
)

func colorizeStr(msg string, c ...color) string {
	colors := ""
	for _, color := range c {
		colors += string(color)
	}
	return fmt.Sprintf("%s%s%s", colors, msg, reset)
}
