package messages

import (
	"fmt"
	"time"
)

// bright purple, cyan, green, yellow
var (
	winTitle  = ColorizeStr("win", Bold, Green, GreenBg)
	noteTitle = ColorizeStr("note", Bold, Cyan, CyanBg)
	fixTitle  = ColorizeStr("fix", Bold, Yellow, YellowBg)

	redTitle          = ColorizeStr("red", Bold, Red, RedBg)
	blackTitle        = ColorizeStr("black", Bold, BlackBg)
	whiteTitle        = ColorizeStr("white", Bold, WhiteBg)
	greenTitle        = ColorizeStr("green", Bold, Green, GreenBg)
	yellowTitle       = ColorizeStr("yellow", Bold, Yellow, YellowBg)
	blueTitle         = ColorizeStr("blue", Bold, Blue, BlueBg)
	purpleTitle       = ColorizeStr("purple", Bold, Purple, PurpleBg)
	cyanTitle         = ColorizeStr("cyan", Bold, Cyan, CyanBg)
	brightBlackTitle  = ColorizeStr("brightBlack", Bold, BrightBlack, BrightBlackBg)
	brightRedTitle    = ColorizeStr("brightRed", Bold, BrightRed, BrightRedBg)
	brightGreenTitle  = ColorizeStr("brightGreen", Bold, BrightGreen, BrightGreenBg)
	brightYellowTitle = ColorizeStr("brightYellow", Bold, BrightYellow, BrightYellowBg)
	brightBlueTitle   = ColorizeStr("brightBlue", Bold, BrightBlue, BrightBlueBg)
	brightPurpleTitle = ColorizeStr("brightPurple", Bold, BrightPurple, BrightPurpleBg)
	brightCyanTitle   = ColorizeStr("brightCyan", Bold, BrightCyan, BrightCyanBg)
	brightWhiteTitle  = ColorizeStr("brightWhite", Bold, BrightWhite, BrightWhiteBg)
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
	tStr = ColorizeStr(tStr, BrightBlack)
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
	date = ColorizeStr(date, BrightPurple)
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
