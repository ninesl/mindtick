package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ninesl/mindtick/messages"
	_ "modernc.org/sqlite"
	//https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc
)

const DBFileName = "store.mindtick"

var COLORDBFILENAME = messages.ColorizeStr(DBFileName, messages.Purple, messages.BrightCyanBg)

func LoadMindtick() (*sql.DB, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get access to directory: %v", err)
	}

	for {
		dbPath := dir + string(os.PathSeparator) + DBFileName
		// fmt.Println("Checking", dbPath)
		if _, err := os.Stat(dbPath); err == nil {
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				return nil, fmt.Errorf("unable to open %s. Is the file corrupted? %v", COLORDBFILENAME, err)
			}
			return db, nil
		}

		parentDir := dir + string(os.PathSeparator) + ".."
		parentDir, err = filepath.Abs(parentDir)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve parent directory: %v", err)
		}
		if parentDir == dir {
			return nil, fmt.Errorf("%s file not found\n%s to create a new mindtick", COLORDBFILENAME, messages.ColorizeStr("mindtick new", messages.BrightGreen))
		}
		dir = parentDir
	}
	// return nil, fmt.Errorf("mindtick file not found. you shouldn't have gotten here")
}

// `mindtick new` command
func New() error {
	// Verify if file exists
	if _, err := os.Stat(DBFileName); os.IsNotExist(err) {
		// do nothing
	} else {
		return fmt.Errorf("%s already exists", COLORDBFILENAME)
	}

	//FIXME: only append to git ignore if DBFileName is not already in there
	// Append to .gitignore if needed
	content, err := os.ReadFile(".gitignore")
	if err == nil && !strings.Contains(string(content), DBFileName) {
		err = os.WriteFile(".gitignore", []byte(string(content)+"\n"+DBFileName+"\n"), 0644)
		if err != nil {
			return fmt.Errorf("failed to update .gitignore: %v", err)
		}
	}

	file, err := os.Create(DBFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s %v", COLORDBFILENAME, err)
	}
	defer file.Close()

	// setup database schema
	db, err := LoadMindtick()
	if err != nil {
		return err
	}
	createSchema(db)

	return fmt.Errorf("%s %s", COLORDBFILENAME, messages.ColorizeStr("intialized", messages.BrightPurple))
}

func Delete() error {
	_, err := LoadMindtick()
	if err != nil {
		return err
	}
	err = os.Remove(DBFileName)
	if err != nil {
		return fmt.Errorf("failed to remove %s\n%v", COLORDBFILENAME, err)
	}
	return fmt.Errorf("%s %s", COLORDBFILENAME, messages.ColorizeStr("deleted", messages.BrightPurple))
}

func createSchema(db *sql.DB) error {
	// schema is
	// messages
	// 	id - serial
	// 	timestamp - time
	// 	msg - string
	// 	msgtype - int enums in messages.go

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		msg TEXT,
		msgtype INT
	)`)
	if err != nil {
		return fmt.Errorf("unable to create mindtick schema: %v", err)
	}
	return nil
}

func AddMessage(db *sql.DB, message messages.Message) error {
	_, err := db.Exec("INSERT INTO messages (timestamp, msg, msgtype) VALUES (?, ?, ?)", message.Timestamp, message.Msg, message.Tag)
	if err != nil {
		return fmt.Errorf("unable to add message: %v", err)
	}
	return nil
}

type Range uint8

const (
	ANYTIME Range = iota
	TODAY
	YESTERDAY
	WEEK
	MONTH
)

var (
	StrToRange = map[string]Range{
		"today":     TODAY,
		"yesterday": YESTERDAY,
		"week":      WEEK,
		"month":     MONTH,
	}
	RangeToStr = map[Range]string{ // dumb
		TODAY:     "today",
		YESTERDAY: "yesterday",
		WEEK:      "week",
		MONTH:     "month",
	}
	RangeToTime = map[Range]func() time.Time{
		TODAY: func() time.Time {
			now := time.Now()
			return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		},
		YESTERDAY: func() time.Time {
			now := time.Now()
			yesterday := now.AddDate(0, 0, -1)
			return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
		},
		WEEK: func() time.Time {
			now := time.Now()
			weekAgo := now.AddDate(0, 0, -7)
			return time.Date(weekAgo.Year(), weekAgo.Month(), weekAgo.Day(), 0, 0, 0, 0, now.Location())
		},
		MONTH: func() time.Time {
			now := time.Now()
			monthAgo := now.AddDate(0, -1, 0)
			return time.Date(monthAgo.Year(), monthAgo.Month(), monthAgo.Day(), 0, 0, 0, 0, now.Location())
		},
	}
	RangeOrder = []Range{TODAY, YESTERDAY, WEEK, MONTH}
)

func Messages(db *sql.DB, tag messages.Tag, rangeType Range) ([]messages.Message, error) {
	var SQLstmt string
	var rows *sql.Rows
	var err error

	if rangeType == ANYTIME && tag == messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages ORDER BY timestamp"
		rows, err = db.Query(SQLstmt)
	}

	if rangeType != ANYTIME && tag == messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages WHERE timestamp > ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, RangeToTime[rangeType]())
	}

	if rangeType == ANYTIME && tag != messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages WHERE msgtype = ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, tag)
	}

	if rangeType != ANYTIME && tag != messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages WHERE msgtype = ? AND timestamp >= ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, tag, RangeToTime[rangeType]())
	}

	if err != nil {
		return nil, fmt.Errorf("unable to query messages: %v", err)
	}
	defer rows.Close()

	return processRows(rows)
}

func processRows(rows *sql.Rows) ([]messages.Message, error) {
	var msgs []messages.Message
	for rows.Next() {
		var msg messages.Message
		err := rows.Scan(&msg.ID, &msg.Timestamp, &msg.Msg, &msg.Tag)
		if err != nil {
			return nil, fmt.Errorf("unable to scan messages: %v", err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// `mindtick edit` command?
func ChangeTimestamp(db *sql.DB, id int, timestamp time.Time) error {
	_, err := db.Exec("UPDATE messages SET timestamp = ? WHERE id = ?", timestamp, id)
	if err != nil {
		return fmt.Errorf("unable to update timestamp: %v", err)
	}
	return nil
}
