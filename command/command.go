package command

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/ninesl/mindtick/messages"
	"github.com/ninesl/mindtick/store"
)

func Exec() {
	if err := processArgs(); err != nil {
		fmt.Println(err)
	}
}

//go:embed version
var version string

var (
	MINDTICK      string = messages.ColorizeStr("mindtick", messages.BrightGreen)
	Ver           string = messages.ColorizeStr(fmt.Sprintf("mindtick %s", version), messages.Bold, messages.BrightRedBg)
	useHelpMsg           = fmt.Sprintf("use %s for more information\n", messages.ColorizeStr("mindtick help", messages.BrightGreen))
	messagePrefix        = "-" // FIXME: how to ignore =:zsh or >:newfilem ' ", other chars, etc? custom prefix?
)

func Version() error {
	fmt.Println(Ver)
	return nil
}
func Help() error {
	var sb strings.Builder
	sb.WriteString(Ver)
	sb.WriteString("\nUsage\n")
	sb.WriteString(messages.ColorizeStr("mindtick command args\n", messages.BrightGreen))
	sb.WriteString("\nCommands\n")
	for _, cmd := range commandOrder {
		sb.WriteString(helpLine(cmd, commandsHelp[cmd]))
	}
	sb.WriteString("\nPlanned Features\n")
	sb.WriteString(plannedFeatureLine("export {tags} {filetype}", "Export all messages to a .pdf/csv/txt file based off specific tags"))
	sb.WriteString(plannedFeatureLine("delete <id>", "Delete a message by id"))
	sb.WriteString(plannedFeatureLine("edit <id> <new message>", "Edit a message by id"))
	sb.WriteString(plannedFeatureLine("{keyword}", "filter by substring"))
	sb.WriteString(plannedFeatureLine("{YYYY-MM-DD}", "filter by date"))

	fmt.Print(sb.String())
	return nil
}

func helpLine(MessageStrategy, description string) string {
	return fmt.Sprintf(
		"%10s\t%s\n",
		messages.ColorizeStr(MessageStrategy, messages.BrightGreen),
		description,
	)
}

func plannedFeatureLine(MessageStrategy, description string) string {
	return fmt.Sprintf(
		"%45s\t%s\n",
		messages.ColorizeStr(MessageStrategy, messages.BrightCyan),
		description,
	)
}

func Ranges() error {
	var sb strings.Builder
	sb.WriteString(helpLine("USAGE: ", messages.ColorizeStr("mindtick view range", messages.BrightPurple)))
	sb.WriteString(helpLine("", messages.ColorizeStr("mindtick view range tag", messages.BrightPurple)))
	for i := range len(store.RangeOrder) {
		date := messages.RenderDate(store.RangeToTime[store.RangeOrder[i]]())
		sb.WriteString(plannedFeatureLine(store.RangeToStr[store.RangeOrder[i]], fmt.Sprintf("filter messages now to %s", date)))
	}
	fmt.Print(sb.String())
	return nil
}

func Tags() error {
	var sb strings.Builder
	sb.WriteString(helpLine("USAGE:", messages.ColorizeStr("mindtick tag -your message", messages.BrightPurple)))
	sb.WriteString(helpLine("", messages.ColorizeStr("mindtick view tag", messages.BrightPurple)))
	sb.WriteString(helpLine("", messages.ColorizeStr("mindtick view tag range", messages.BrightPurple)))
	for _, tag := range messages.TagOrder {
		sb.WriteString(messages.Tags[tag] + " ")
	}
	fmt.Println(sb.String())
	return nil
}

var (
	commands = map[string]func() error{
		"help":    Help,
		"version": Version,
		"new":     store.New,
		"delete":  store.Delete,
		"tag":     AddMessage,
		"view":    View,
		"tags":    Tags,
		"ranges":  Ranges,
	}
	commandsHelp = map[string]string{
		"help":    "Display this help message",
		"version": fmt.Sprintf("Display the current version of %s", MINDTICK),
		"new":     fmt.Sprintf("Create a new %s file in the current directory", store.COLORDBFILENAME),
		"delete":  fmt.Sprintf("Delete the %s file in the current directory", store.COLORDBFILENAME),
		"tag":     fmt.Sprintf("%s | adds a message", messages.ColorizeStr("-your message", messages.BrightPurple)),
		"view":    fmt.Sprintf("optional: %s | Display messages by tag and/or range", messages.ColorizeStr("tag range", messages.BrightPurple)),
		"tags":    fmt.Sprintf("Display all available tags, used in %s and %s", messages.ColorizeStr("view", messages.BrightGreen), messages.ColorizeStr("tag", messages.BrightGreen)),
		"ranges":  "Display all available ranges",
	}
	commandOrder = []string{"version", "help", "new", "delete", "tag", "view", "tags", "ranges"}
)

func processArgs() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("mindtick requires at least one argument, %s", useHelpMsg)
	}

	if _, ok := messages.StrToTag[strings.ToLower(os.Args[1])]; ok { // case insensitivity for tags
		return commands["tag"]()
	}

	if len(os.Args) > 1 {
		if cmd, ok := commands[os.Args[1]]; ok {
			return cmd()
		}
	} else {
		return fmt.Errorf("mindtick requires at least one argument, %s", useHelpMsg)
	}

	return fmt.Errorf("unknown mindtick argument %s, %s", messages.ColorizeStr(strings.Join(os.Args[1:], " "), messages.BrightPurple), useHelpMsg)
}

func View() error {
	args := os.Args
	size := len(args)

	if size > 4 {
		return fmt.Errorf("too many arguments for view, %s", useHelpMsg)
	}

	db, err := store.LoadMindtick()
	if err != nil {
		return err
	}

	if size == 2 { // default behavior
		msgs, err := store.Messages(db, messages.ANYTAG, store.ANYTIME)
		if err != nil {
			return err
		}

		if len(msgs) == 0 {
			return fmt.Errorf("%s is empty, %s", messages.ColorizeStr(store.DBFileName, messages.BrightRed), useHelpMsg)
		}
		messages.RenderMessages(msgs...)
		return nil
	}

	rangeType := store.StrToRange[args[2]]
	msgType := messages.StrToTag[args[2]]

	if rangeType == store.ANYTIME && msgType == messages.ANYTAG {
		return fmt.Errorf("unknown view argument %s, %s", messages.ColorizeStr(args[2], messages.BrightPurple), useHelpMsg)
	}

	if size == 4 {
		// if range is not found, check if it's a tag
		if rangeType == store.ANYTIME {
			rangeType = store.StrToRange[args[3]]
			if rangeType == store.ANYTIME {
				var ranges []string
				for k := range store.StrToRange {
					ranges = append(ranges, k)
				}
				return fmt.Errorf("unknown view range %s\nvalid ranges are %v", messages.ColorizeStr(args[3], messages.BrightPurple), messages.ColorizeStr(strings.Join(ranges, ", "), messages.BrightGreen))
			}
		} else if msgType == messages.ANYTAG {
			msgType = messages.StrToTag[args[3]]
			if msgType == messages.ANYTAG {
				var msgTypes []string
				for k := range messages.StrToTag {
					msgTypes = append(msgTypes, k)
				}
				return fmt.Errorf("unknown view tag %s\nvalid tags are %v", messages.ColorizeStr(args[3], messages.BrightPurple), messages.ColorizeStr(strings.Join(msgTypes, ", "), messages.BrightGreen))
			}
		}
	}

	msgs, err := store.Messages(db, msgType, rangeType)
	if err != nil {
		return err
	}

	if len(msgs) == 0 {
		return fmt.Errorf("no messages found with %s", messages.ColorizeStr(strings.Join(os.Args[2:], " "), messages.BrightPurple))
	}
	messages.RenderMessages(msgs...)
	return nil
}

func AddMessage() error {
	tagCmd := strings.ToLower(os.Args[1]) // case insensitivity for tags
	db, err := store.LoadMindtick()
	if err != nil {
		return err
	}
	if len(os.Args) < 3 {
		return fmt.Errorf("mindtick %s must have a message, %s", messages.ColorizeStr(tagCmd, messages.BrightPurple), useHelpMsg)
	}
	if os.Args[2][0:1] != messagePrefix {
		tip := fmt.Sprintf("mindtick %s %v%s", tagCmd, messagePrefix, strings.Join(os.Args[2:], " "))
		return fmt.Errorf("mindtick messages must start with %v\nexample usage: %s\n%s", messages.ColorizeStr(messagePrefix, messages.BrightGreen), messages.ColorizeStr(tip, messages.BrightGreen), useHelpMsg)
	}

	// fmt.Println(os.Args)
	os.Args[2] = strings.Replace(os.Args[2], messagePrefix, "", 1)

	// concat all arguments after the msgType
	var argMsgs []string
	for i := 2; i < len(os.Args); i++ {
		argMsgs = append(argMsgs, os.Args[i])
	}

	var argMsg = strings.Join(argMsgs, " ")
	msg, err := messages.NewMessage(tagCmd, argMsg)
	if err != nil {
		return fmt.Errorf("%s, %s", messages.ColorizeStr(err.Error(), messages.BrightRed), useHelpMsg)
	}

	err = store.AddMessage(db, msg)
	if err != nil {
		return fmt.Errorf("%s, %s", messages.ColorizeStr(err.Error(), messages.BrightRed), useHelpMsg)
	}

	//check to see if message was added
	//	get the id, check if id exists, if true
	//	print message (msg not loading the whole msg from the DB)
	// else
	// 	return error
	fmt.Println(messages.RenderMsg(msg, false))
	return nil
}
