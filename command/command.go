package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/ninesl/mindtick/messages"
	"github.com/ninesl/mindtick/store"
)

var (
	useHelpMsg    = fmt.Sprintf("use %s for more information\n", messages.ColorizeStr("mindtick help", messages.BrightGreen))
	messagePrefix = '>'
)

func helpLine(command, description string) string {
	return fmt.Sprintf(
		"%10s\t%s\n",
		messages.ColorizeStr(command, messages.BrightGreen),
		description,
	)
}

func help() {
	var sb strings.Builder
	sb.WriteString("Usage\n")
	sb.WriteString(messages.ColorizeStr("mindtick <command>\n", messages.BrightGreen))
	sb.WriteString("\nCommands\n")
	sb.WriteString(helpLine("new", "Create a new mindtick file in the current directory"))
	sb.WriteString(helpLine("delete", "Delete the mindtick file in the current directory"))
	sb.WriteString(helpLine("help", "Display this help message"))

	fmt.Print(sb.String())
}

func ProcessArgs() {
	if len(os.Args) < 2 {
		fmt.Printf("No arguments provided, %s", useHelpMsg)
		return
	}

	switch os.Args[1] {
	case "new":
		err := store.New()
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		}
		fmt.Println(messages.ColorizeStr("mindtick intialized", messages.BrightPurple))
	case "delete":
		err := store.Delete()
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		} else {
			fmt.Println(messages.ColorizeStr("mindtick deleted", messages.BrightPurple))
		}
	case "help":
		help()

	// list of note commands is hard to maintain and extend with current implementation
	case "win": // new note, usage is like `mindtick {msgType} "note content"`
	case "task":
	case "note":
	case "fix":
		processMessage()
	default:
		fmt.Printf("unknown mindtick argument %s, %s", messages.ColorizeStr(os.Args[1], messages.BrightPurple), useHelpMsg)
	}
}

func processMessage() {
	if len(os.Args) < 3 {
		fmt.Printf("mindtick %s must have a message, %s", messages.ColorizeStr(os.Args[1], messages.BrightPurple), useHelpMsg)
		return
	}
	if rune(os.Args[2][0]) != messagePrefix {
		tip := fmt.Sprintf("mindtick %s >{your message here}", os.Args[1])
		fmt.Printf("mindtick messages must start with >. example usage: %s\n%s", messages.ColorizeStr(tip, messages.BrightGreen), useHelpMsg)
		return
	}
	// concat all arguments after the msgType
	var argMsgs []string
	for i := 3; i < len(os.Args); i++ {
		argMsgs = append(argMsgs, os.Args[i])
	}
	var argMsg = strings.Join(argMsgs, " ")
	argMsg = strings.Replace(argMsg, ">", "", 1) // remove the > from the beginning
	msg, err := messages.NewMessage(os.Args[1], argMsg)
	if err != nil {
		fmt.Printf("%s, %s", messages.ColorizeStr(err.Error(), messages.BrightRed), useHelpMsg)
		return
	}
	err = store.AddMessage(msg)
	if err != nil {
		fmt.Printf("%s, %s", messages.ColorizeStr(err.Error(), messages.BrightRed), useHelpMsg)
		return
	}
}
