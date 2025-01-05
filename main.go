package main

import (
	"github.com/ninesl/mindtick/command"
)

func main() {
	// msgs := []messages.Message{
	// 	{Msg: "This is a win message", MsgType: messages.WIN},
	// 	{Msg: "This is a note message", MsgType: messages.NOTE},
	// 	{Msg: "This is a fix message", MsgType: messages.FIX},
	// }
	// for i := range msgs {
	// 	msgs[i].Timestamp = time.Now().Add(time.Duration(rand.Intn(1000)) * time.Hour)
	// }
	// messages.RenderMessages(msgs...)

	command.ProcessArgs()

	// messages.PrintAllTitles()
}
