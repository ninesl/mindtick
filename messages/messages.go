package messages

import (
	"fmt"
	"time"
)

// bright purple, cyan, green, yellow
var (
	// winTag  = ColorizeStr("  win", Bold, BrightWhite, GreenBg)
	// winBg   = ColorizeStr("     ", Bold, BrightWhite, GreenBg)
	// noteTag = ColorizeStr(" note", Bold, BrightWhite, CyanBg)
	// noteBg  = ColorizeStr("     ", Bold, BrightWhite, CyanBg)
	// fixTag  = ColorizeStr("  fix", Bold, BrightWhite, YellowBg)
	// fixBg   = ColorizeStr("     ", Bold, BrightWhite, YellowBg)
	// taskTag = ColorizeStr(" task", Bold, BrightWhite, PurpleBg)
	// taskBg  = ColorizeStr("     ", Bold, BrightWhite, PurpleBg)

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
	URL
	WORK
)

// TODO: FIXME: refactor this mess when impl custom tags
// To add a hardcoded tag, follow these steps:
// 1. Add a new constant to the Tag type
// 2. Add a new constant to the Tags map
// 3. Add a new constant to the bgs map
// 4. Add a new constant to the StrToTag map
var (
	winTag   = ColorizeStr("  win", BrightGreenBg, Bold, White)
	winBg    = ColorizeStr("     ", BrightGreenBg, Bold, White)
	noteTag  = ColorizeStr(" note", CyanBg, Bold, White)
	noteBg   = ColorizeStr("     ", CyanBg, Bold, White)
	fixTag   = ColorizeStr("  fix", BrightYellowBg, Bold, White)
	fixBg    = ColorizeStr("     ", BrightYellowBg, Bold, White)
	taskTag  = ColorizeStr(" task", BrightPurpleBg, Bold, White)
	taskBg   = ColorizeStr("     ", BrightPurpleBg, Bold, White)
	urlTag   = ColorizeStr("  url", BlackBg, Bold, Blue)
	urlBg    = ColorizeStr("     ", BlackBg, Bold, Blue)
	workTag  = ColorizeStr(" work", BrightWhiteBg, Bold, Black)
	workBg   = ColorizeStr("     ", BrightWhiteBg, Bold, Black)
	TagOrder = []Tag{WIN, NOTE, FIX, TASK, URL, WORK}
	StrToTag = map[string]Tag{
		"win":  WIN,
		"note": NOTE,
		"fix":  FIX,
		"task": TASK,
		"url":  URL,
		"work": WORK,
	}
	Tags = map[Tag]string{
		WIN:  winTag,
		NOTE: noteTag,
		FIX:  fixTag,
		TASK: taskTag,
		URL:  urlTag,
		WORK: workTag,
	}
	bgs = map[Tag]string{
		WIN:  winBg,
		NOTE: noteBg,
		FIX:  fixBg,
		TASK: taskBg,
		URL:  urlBg,
		WORK: workBg,
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

func NewMessage(tagStr string, msg string) (Message, error) {
	tag := StrToTag[tagStr]

	// if tag == ANYTAG { // should never get here
	// 	panic("unknown view tag " + tagStr)
	// 	// return Message{}, fmt.Errorf("unknown view tag %s", tagStr)
	// }

	return Message{
		Timestamp: time.Now(),
		Msg:       msg,
		Tag:       tag,
	}, nil
}
