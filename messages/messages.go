package messages

import (
	"fmt"
	"time"
)

// bright purple, cyan, green, yellow
var (
	winTag  = ColorizeStr("  win", Bold, Yellow, GreenBg)
	winBg   = ColorizeStr("     ", Bold, Yellow, GreenBg)
	noteTag = ColorizeStr(" note", Bold, Yellow, CyanBg)
	noteBg  = ColorizeStr("     ", Bold, Yellow, CyanBg)
	fixTag  = ColorizeStr("  fix", Bold, Yellow, YellowBg)
	fixBg   = ColorizeStr("     ", Bold, Yellow, YellowBg)
	taskTag = ColorizeStr(" task", Bold, Yellow, PurpleBg)
	taskBg  = ColorizeStr("     ", Bold, Yellow, PurpleBg)

	redTag          = ColorizeStr("red", Bold, Red, RedBg)
	blackTag        = ColorizeStr("black", Bold, BlackBg)
	whiteTag        = ColorizeStr("white", Bold, WhiteBg)
	greenTag        = ColorizeStr("green", Bold, Green, GreenBg)
	yellowTag       = ColorizeStr("yellow", Bold, Yellow, YellowBg)
	blueTag         = ColorizeStr("blue", Bold, Blue, BlueBg)
	purpleTag       = ColorizeStr("purple", Bold, Purple, PurpleBg)
	cyanTag         = ColorizeStr("cyan", Bold, Cyan, CyanBg)
	brightBlackTag  = ColorizeStr("brightBlack", Bold, BrightBlack, BrightBlackBg)
	brightRedTag    = ColorizeStr("brightRed", Bold, BrightRed, BrightRedBg)
	brightGreenTag  = ColorizeStr("brightGreen", Bold, BrightGreen, BrightGreenBg)
	brightYellowTag = ColorizeStr("brightYellow", Bold, BrightYellow, BrightYellowBg)
	brightBlueTag   = ColorizeStr("brightBlue", Bold, BrightBlue, BrightBlueBg)
	brightPurpleTag = ColorizeStr("brightPurple", Bold, BrightPurple, BrightPurpleBg)
	brightCyanTag   = ColorizeStr("brightCyan", Bold, BrightCyan, BrightCyanBg)
	brightWhiteTag  = ColorizeStr("brightWhite", Bold, BrightWhite, BrightWhiteBg)

	Tags = map[Tag]string{
		WIN:  winTag,
		NOTE: noteTag,
		FIX:  fixTag,
		TASK: taskTag,
	}
	bgs = map[Tag]string{
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

// Prints all tags for display testing purposes.
func PrintAllTags() {
	tags := []string{
		winTag, noteTag, fixTag, redTag, blackTag, whiteTag, greenTag, yellowTag,
		blueTag, purpleTag, cyanTag, brightBlackTag, brightRedTag, brightGreenTag,
		brightYellowTag, brightBlueTag, brightPurpleTag, brightCyanTag, brightWhiteTag,
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}

type Tag uint8

const (
	ANYTAG Tag = iota
	WIN
	NOTE
	FIX
	TASK
)

var (
	StrToTag = map[string]Tag{
		"win":  WIN,
		"note": NOTE,
		"fix":  FIX,
		"task": TASK,
	}
)

// id INTEGER PRIMARY KEY AUTOINCREMENT,
// timestamp DATETIME,
// msg TEXT,
// msgtype INT
type Message struct {
	Timestamp time.Time `db:"timestamp"`
	Msg       string    `db:"msg"`
	ID        int       `db:"id"`
	Tag       Tag       `db:"msgtype"`
}

func renderTime(t time.Time) string {
	time := t.Format("03:04 PM")

	tStr := fmt.Sprintf("%8s", time)
	tStr = ColorizeStr(tStr, BrightBlack)
	return tStr
}

func RenderTag(msgType Tag, bgOnly bool) string {
	var tag string

	if bgOnly {
		tag = bgs[msgType]
	} else {
		tag = Tags[msgType]
	}
	tag = fmt.Sprintf("%*s", 23, tag)
	// 23 is the length of the longest tag bc of ANSI color

	return tag
}

func RenderMsg(msg Message, bgOnly bool) string {
	var (
		tag  = RenderTag(msg.Tag, bgOnly)
		time = renderTime(msg.Timestamp)
	)

	return fmt.Sprintf("%s %s     %s", tag, time, msg.Msg)
}

func RenderDate(d time.Time) string {
	date := fmt.Sprintf("[ %s ]", d.Format("Jan 02, 2006"))
	date = ColorizeStr(date, BrightPurple)
	return date
}

// will always be sorted by timestamp
func RenderMessages(msgs ...Message) {
	curDate := msgs[0].Timestamp
	curType := ANYTAG
	fmt.Println(RenderDate(curDate))

	for i := range msgs {
		if curDate.Day() != msgs[i].Timestamp.Day() {
			curDate = msgs[i].Timestamp
			curType = ANYTAG
			fmt.Println("\n" + RenderDate(curDate))
		}

		if curType != msgs[i].Tag {
			curType = msgs[i].Tag
			fmt.Println(RenderMsg(msgs[i], BGTITLE))
		} else {
			fmt.Println(RenderMsg(msgs[i], ONLYBG))
		}
	}
}

func NewMessage(msgTypeStr string, msg string) (Message, error) {
	var msgType Tag

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
		return Message{}, fmt.Errorf("invalid message type: %s", ColorizeStr(msgTypeStr, White, BrightPurpleBg))
	}

	return Message{
		Timestamp: time.Now(),
		Msg:       msg,
		Tag:       msgType,
	}, nil
}
