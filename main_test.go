// main_test.go
package main

import (
	"math/rand/v2"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/ninesl/mindtick/messages"
)

func TestRenderOutput(t *testing.T) {
	t.Run("generates and renders random messages", func(t *testing.T) {
		msgs := generateTestMessages(15)

		// These don't need to be tested, but we will use it later for store/sqllite.go testing

		// Verify message count
		if len(msgs) != 15 {
			t.Errorf("expected 15 messages, got %d", len(msgs))
		}

		// Verify timestamps are properly ordered
		for i := 1; i < len(msgs); i++ {
			if msgs[i].Timestamp.Before(msgs[i-1].Timestamp) {
				t.Errorf("messages not properly sorted by timestamp at index %d", i)
			}
		}

		messages.RenderMessages(msgs...) // This will print to terminal

		// Verify all message types are present
		types := make(map[messages.MessageType]bool)
		for _, msg := range msgs {
			types[msg.MsgType] = true
		}
		expectedTypes := []messages.MessageType{messages.WIN, messages.NOTE, messages.FIX, messages.TASK}
		for _, expectedType := range expectedTypes {
			if !types[expectedType] {
				t.Errorf("missing message type: %v", expectedType)
			}
		}
	})
}

func generateTestMessages(count int) []messages.Message {
	msgs := []messages.Message{}
	msgTypes := []string{"win", "note", "fix", "task"}
	msgCount := map[string]int{"win": 0, "note": 0, "fix": 0, "task": 0}

	for i := 0; i < count; i++ {
		msgType := msgTypes[rand.IntN(len(msgTypes))]
		var msg messages.Message
		switch msgType {
		case "win":
			msg, _ = messages.NewMessage("win", "win message "+strconv.Itoa(msgCount[msgType]))
		case "note":
			msg, _ = messages.NewMessage("note", "note message "+strconv.Itoa(msgCount[msgType]))
		case "fix":
			msg, _ = messages.NewMessage("fix", "fix message "+strconv.Itoa(msgCount[msgType]))
		case "task":
			msg, _ = messages.NewMessage("task", "task message "+strconv.Itoa(msgCount[msgType]))
		}
		msgs = append(msgs, msg)
		msgCount[msgType]++
	}

	rand.Shuffle(len(msgs), func(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] })

	for i := range msgs {
		randomTime := time.Now().Add(-time.Duration(rand.IntN(48)) * time.Hour)
		msgs[i].Timestamp = randomTime
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp.Before(msgs[j].Timestamp)
	})

	return msgs
}
