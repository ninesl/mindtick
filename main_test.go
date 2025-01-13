// main_test.go
package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ninesl/mindtick/command"
	"github.com/ninesl/mindtick/messages"
	"github.com/ninesl/mindtick/store"
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
		types := make(map[messages.Tag]bool)
		for _, msg := range msgs {
			types[msg.Tag] = true
		}
		expectedTypes := []messages.Tag{messages.WIN, messages.NOTE, messages.FIX, messages.TASK}
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
			msg, _ = messages.NewMessage("win", "win message")
		case "note":
			msg, _ = messages.NewMessage("note", "note message")
		case "fix":
			msg, _ = messages.NewMessage("fix", "fix message")
		case "task":
			msg, _ = messages.NewMessage("task", "task message")
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
func addArgs(args []string) {
	os.Args = append(os.Args, args...)
}

var (
	tagMessages = map[string][]string{
		"win": {
			"-deployed new feature to production",
			"-achieved 100% test coverage",
			"-completed user authentication system",
			"-successful client demo",
			"-optimized database queries by 50%",
			"-merged major feature branch",
			"-secured new enterprise client",
			"-reduced API latency by 40%",
			"-implemented zero-downtime deployments",
			"-completed security audit",
			"-achieved AWS certification",
			"-launched mobile app v2.0",
			"-migrated legacy system successfully",
			"-scaled system to 1M users",
			"-won hackathon first place",
		},
		"note": {
			"-need to update documentation",
			"-api rate limits changed to 1000/min",
			"-team meeting moved to Thursdays",
			"-switching to new logging framework",
			"-considering kubernetes migration",
			"-database backup scheduled weekly",
			"-new API version available",
			"-cloud costs increased 15%",
			"-team expanding next month",
			"-planning system architecture v2",
			"-evaluating new cache solution",
			"-documentation needs review",
			"-monitoring system upgrade needed",
			"-customer feedback session tomorrow",
			"-tech debt assessment completed",
		},
		"fix": {
			"-resolved memory leak in worker pool",
			"-fixed broken CI pipeline",
			"-patched security vulnerability",
			"-corrected timezone handling bug",
			"-fixed race condition in cache layer",
			"-resolved null pointer exception",
			"-eliminated database deadlock",
			"-fixed authentication bypass",
			"-resolved session handling issue",
			"-fixed data corruption bug",
			"-patched XSS vulnerability",
			"-fixed broken deployment script",
			"-resolved API versioning conflict",
			"-fixed memory overflow issue",
			"-corrected data validation logic",
		},
		"task": {
			"-implement rate limiting",
			"-update dependencies",
			"-add error monitoring",
			"-setup staging environment",
			"-create backup strategy",
			"-write integration tests",
			"-implement CI/CD pipeline",
			"-configure load balancer",
			"-setup monitoring alerts",
			"-implement search feature",
			"-add user analytics",
			"-create api documentation",
			"-implement caching layer",
			"-setup disaster recovery",
			"-add performance metrics",
		},
	}
)

func mindtick(input string) {
	args := strings.Fields(input)
	os.Args = append([]string{"mindtick"}, args...)
	fmt.Println()
	fmt.Println("Running command:")
	fmt.Println(messages.ColorizeStr("mindtick "+strings.Join(args, " "), messages.Cyan))
	fmt.Println()
	command.Exec()
}

func TestRenderWithCustomArgs(t *testing.T) {
	// Save original
	oldArgs := os.Args

	var args []string

	mindtick("new")
	mindtick("tags")
	mindtick("ranges")

	db, err := store.LoadMindtick()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ids := 100
	tags := []string{"win", "note", "fix", "task"}
	// add ids random messages
	for id := 0; id < ids; id++ {
		os.Args = []string{"mindtick"}
		tag := tags[rand.IntN(4)]
		args = append(args, tag, tagMessages[tag][rand.IntN(len(tagMessages[tag]))])

		addArgs(args)
		command.Exec()
		args = args[:0]

	}

	for id := range ids - 1 {
		behindNow := time.Now().Add(-time.Duration(rand.IntN(7*24)) * time.Hour).Add(-time.Duration(rand.IntN(60)) * time.Minute)
		store.ChangeTimestamp(db, id, behindNow)
	}

	fmt.Println()
	fmt.Println("Running command:")
	fmt.Println(messages.ColorizeStr(fmt.Sprintf("mindtick %v -%d random messages", tags, ids), messages.Cyan))

	// daysAgo := rand.IntN(15) + 15
	// mindtick(fmt.Sprintf("win -%d days ago!", daysAgo))
	// store.ChangeTimestamp(db, ids, time.Now().Add(-((time.Hour * 24) * time.Duration(daysAgo))))
	// mindtick("win -Hello World!") //FIXME: ids not getting right message
	// store.ChangeTimestamp(db, ids+1, time.Now().Add(-time.Hour*24*time.Duration(30)))

	mindtick("view")
	// mindtick("view task")
	mindtick("view yesterday")
	mindtick("view yesterday task")

	mindtick("help")

	mindtick("delete")
	// IMPORTANT: Restore original state
	defer func() { os.Args = oldArgs }()

}
