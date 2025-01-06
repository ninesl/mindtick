package messages

import (
	"fmt"
	"time"
)

// bright purple, cyan, green, yellow
var (
	winTitle  = ColorizeStr("  win", Bold, Yellow, GreenBg)
	winBg     = ColorizeStr("     ", Bold, Yellow, GreenBg)
	noteTitle = ColorizeStr(" note", Bold, Yellow, CyanBg)
	noteBg    = ColorizeStr("     ", Bold, Yellow, CyanBg)
	fixTitle  = ColorizeStr("  fix", Bold, Yellow, YellowBg)
	fixBg     = ColorizeStr("     ", Bold, Yellow, YellowBg)
	taskTitle = ColorizeStr(" task", Bold, Yellow, PurpleBg)
	taskBg    = ColorizeStr("     ", Bold, Yellow, PurpleBg)

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

	titles = map[MessageType]string{
		WIN:  winTitle,
		NOTE: noteTitle,
		FIX:  fixTitle,
		TASK: taskTitle,
	}
	bgs = map[MessageType]string{
		WIN:  winBg,
		NOTE: noteBg,
		FIX:  fixBg,
		TASK: taskBg,
	}
)

const (
	BGTITLE = false
	ONLYBG  = true
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
	NONE MessageType = iota
	WIN
	NOTE
	FIX
	TASK
)

// id INTEGER PRIMARY KEY AUTOINCREMENT,
// timestamp DATETIME,
// msg TEXT,
// msgtype INT
type Message struct {
	Timestamp time.Time   `db:"timestamp"`
	Msg       string      `db:"msg"`
	ID        int         `db:"id"`
	MsgType   MessageType `db:"msgtype"`
}

func renderTime(t time.Time) string {
	time := t.Format("03:04 PM")

	tStr := fmt.Sprintf("%8s", time)
	tStr = ColorizeStr(tStr, BrightBlack)
	return tStr
}

func RenderTitle(msgType MessageType, bgOnly bool) string {
	var title string

	if bgOnly {
		title = bgs[msgType]
	} else {
		title = titles[msgType]
	}
	title = fmt.Sprintf("%*s", 23, title)
	// 23 is the length of the longest title bc of ANSI color

	return title
}

func RenderMsg(msg Message, bgOnly bool) string {
	var (
		title = RenderTitle(msg.MsgType, bgOnly)
		time  = renderTime(msg.Timestamp)
	)

	return fmt.Sprintf("%s %s     %s", title, time, msg.Msg)
}

func renderDate(d time.Time) string {
	date := fmt.Sprintf("[ %s ]", d.Format("Jan 02, 2006"))
	date = ColorizeStr(date, BrightPurple)
	return date
}

// will always be sorted by timestamp
func RenderMessages(msgs ...Message) {
	curDate := msgs[0].Timestamp
	curType := NONE
	fmt.Println(renderDate(curDate))

	for i := range msgs {
		if curDate.Day() != msgs[i].Timestamp.Day() {
			curDate = msgs[i].Timestamp
			curType = NONE
			fmt.Println("\n" + renderDate(curDate))
		}

		if curType != msgs[i].MsgType {
			curType = msgs[i].MsgType
			fmt.Println(RenderMsg(msgs[i], BGTITLE))
		} else {
			fmt.Println(RenderMsg(msgs[i], ONLYBG))
		}
	}
}

func NewMessage(msgTypeStr string, msg string) (Message, error) {
	var msgType MessageType

	switch msgTypeStr {
	case "win":
		msgType = WIN
	case "note":
		msgType = NOTE
	case "fix":
		msgType = FIX
	case "task":
		msgType = TASK
	default:
		return Message{}, fmt.Errorf("invalid message type")
	}

	return Message{
		Timestamp: time.Now(),
		Msg:       msg,
		MsgType:   msgType,
	}, nil
}
