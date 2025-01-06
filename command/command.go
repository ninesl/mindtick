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
	messagePrefix = "-" // how to ignore =:zsh or >:newfile, etc? custom prefix?
)

func helpLine(command, description string) string {
	return fmt.Sprintf(
		"%10s\t%s\n",
		messages.ColorizeStr(command, messages.BrightGreen),
		description,
	)
}

/*
unsure if this true/false is a good way to handle errors

	returns `true` if there was an error, prints the error
	returns `false` if there was no error
*/
func isHandleGenericError(err error) bool {
	if err != nil {
		fmt.Print(messages.ColorizeStr(err.Error(), messages.BrightRed))
		return true
	}
	return false
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

	// fmt.Println(os.Args[1])

	switch os.Args[1] {
	case "new":
		if len(os.Args) > 2 {
			fmt.Printf("%s does not take any arguments, %s", messages.ColorizeStr("mindtick new", messages.BrightGreen), useHelpMsg)
		}
		err := store.New()
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		}
		fmt.Println(messages.ColorizeStr("mindtick intialized", messages.BrightPurple))
	case "delete":
		if len(os.Args) > 2 {
			fmt.Printf("%s does not take any arguments, %s", messages.ColorizeStr("mindtick delete", messages.BrightGreen), useHelpMsg)
		}
		err := store.Delete()
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		} else {
			fmt.Println(messages.ColorizeStr("mindtick deleted", messages.BrightPurple))
		}
	case "help":
		help()

	case "win", "task", "note", "fix": // combine all message types into a single case
		err := processMessage()
		if !isHandleGenericError(err) {
			return
		}
	case "view":
		db, err := store.LoadMindtick()
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		}
		msgs, err := store.GetMessages(db) // TODO: view by date, type, etc
		if err != nil {
			fmt.Println(messages.ColorizeStr(err.Error(), messages.BrightRed))
			return
		}
		messages.RenderMessages(msgs...)
	default:
		fmt.Printf("unknown mindtick argument %s, %s", messages.ColorizeStr(os.Args[1], messages.BrightPurple), useHelpMsg)
	}
}

func processMessage() error {
	db, err := store.LoadMindtick()
	if err != nil {
		return err
	}

	if len(os.Args) < 3 {
		return fmt.Errorf("mindtick %s must have a message, %s", messages.ColorizeStr(os.Args[1], messages.BrightPurple), useHelpMsg)
	}
	if os.Args[2][0:1] != messagePrefix {
		tip := fmt.Sprintf("mindtick %s %v{your message here}", os.Args[1], messagePrefix)
		return fmt.Errorf("mindtick messages must start with %v. example usage: %s\n%s", messagePrefix, messages.ColorizeStr(tip, messages.BrightGreen), useHelpMsg)
	}

	// fmt.Println(os.Args)
	os.Args[2] = strings.Replace(os.Args[2], messagePrefix, "", 1)

	// concat all arguments after the msgType
	var argMsgs []string
	for i := 2; i < len(os.Args); i++ {
		argMsgs = append(argMsgs, os.Args[i])
	}

	var argMsg = strings.Join(argMsgs, " ")
	argMsg = strings.Replace(argMsg, ">", "", 1) // remove the > from the beginning
	msg, err := messages.NewMessage(os.Args[1], argMsg)
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
