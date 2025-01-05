package messages

import (
	"fmt"
	"time"
)

// bright purple, cyan, green, yellow
var (
	winTitle  = colorizeStr("win", bold, green, greenBg)
	noteTitle = colorizeStr("note", bold, cyan, cyanBg)
	fixTitle  = colorizeStr("fix", bold, yellow, yellowBg)

	redTitle          = colorizeStr("red", bold, red, redBg)
	blackTitle        = colorizeStr("black", bold, blackBg)
	whiteTitle        = colorizeStr("white", bold, whiteBg)
	greenTitle        = colorizeStr("green", bold, green, greenBg)
	yellowTitle       = colorizeStr("yellow", bold, yellow, yellowBg)
	blueTitle         = colorizeStr("blue", bold, blue, blueBg)
	purpleTitle       = colorizeStr("purple", bold, purple, purpleBg)
	cyanTitle         = colorizeStr("cyan", bold, cyan, cyanBg)
	brightBlackTitle  = colorizeStr("brightBlack", bold, brightBlack, brightBlackBg)
	brightRedTitle    = colorizeStr("brightRed", bold, brightRed, brightRedBg)
	brightGreenTitle  = colorizeStr("brightGreen", bold, brightGreen, brightGreenBg)
	brightYellowTitle = colorizeStr("brightYellow", bold, brightYellow, brightYellowBg)
	brightBlueTitle   = colorizeStr("brightBlue", bold, brightBlue, brightBlueBg)
	brightPurpleTitle = colorizeStr("brightPurple", bold, brightPurple, brightPurpleBg)
	brightCyanTitle   = colorizeStr("brightCyan", bold, brightCyan, brightCyanBg)
	brightWhiteTitle  = colorizeStr("brightWhite", bold, brightWhite, brightWhiteBg)
)

func PrintAllTitles() {
	titles := []string{
		winTitle, noteTitle, fixTitle, redTitle, blackTitle, whiteTitle, greenTitle, yellowTitle,
		blueTitle, purpleTitle, cyanTitle, brightBlackTitle, brightRedTitle, brightGreenTitle,
		brightYellowTitle, brightBlueTitle, brightPurpleTitle, brightCyanTitle, brightWhiteTitle,
	}

	for _, title := range titles {
		fmt.Println(title)
	}
}

type MessageType uint8

const (
	WIN MessageType = iota
	NOTE
	FIX
)

type Message struct {
	Timestamp time.Time
	Msg       string
	ID        int
	MsgType   MessageType
}

func renderTime(t time.Time) string {
	time := t.Format("03:04 PM")

	tStr := fmt.Sprintf("[ %8s ]", time)
	tStr = colorizeStr(tStr, brightBlack)
	return tStr
}

func renderTitle(msgType MessageType) string {
	var title string

	switch msgType {
	case WIN:
		title = winTitle
	case NOTE:
		title = noteTitle
	case FIX:
		title = fixTitle
	default:
		title = ""
	}
	title = fmt.Sprintf("%*s", 23, title)
	// 23 is the length of the longest title bc of ANSI color

	return title
}

func renderMsg(msg Message) string {
	var (
		title = renderTitle(msg.MsgType)
		time  = renderTime(msg.Timestamp)
	)

	return fmt.Sprintf("%s %s %s %d", title, time, msg.Msg, len(title))
}

func renderDate(d time.Time) string {
	date := fmt.Sprintf("[ %s ]", d.Format("Jan 02, 2006"))
	date = colorizeStr(date, brightPurple)
	return date
}

// will always be sorted by timestamp
func RenderMessages(msgs ...Message) {
	curDate := msgs[0].Timestamp
	fmt.Println(renderDate(curDate))

	for i := range msgs {
		if curDate.Day() != msgs[i].Timestamp.Day() {
			curDate = msgs[i].Timestamp
			fmt.Println("\n" + renderDate(curDate))
		}
		fmt.Println(renderMsg(msgs[i]))
	}
}
