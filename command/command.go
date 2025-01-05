package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/ninesl/mindtick/messages"
	"github.com/ninesl/mindtick/store"
)

var useHelpMsg = fmt.Sprintf("use %s for more information\n", messages.ColorizeStr("mindtick help", messages.BrightGreen))

func helpLine(command, description string) string {
	return fmt.Sprintf(
		"%s%s\t%s\n",
		messages.ColorizeStr("  ", messages.BrightGreen),
		messages.ColorizeStr(command, messages.BrightGreen),
		description,
	)
}

func help() {
	var sb strings.Builder
	sb.WriteString("Usage: ")
	sb.WriteString(messages.ColorizeStr("mindtick <command>\n", messages.BrightGreen))
	sb.WriteString("Commands:\n")
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
		// err := store.Delete()
	case "help":
		help()
	default:
		fmt.Printf("unknown mindtick argument %s, %s", messages.ColorizeStr(os.Args[1], messages.BrightPurple), useHelpMsg)
	}
}
