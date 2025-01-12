package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ninesl/mindtick/messages"
	_ "modernc.org/sqlite"
	//https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc
)

const dbFileName = "store.mindtick"

func LoadMindtick() (*sql.DB, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get access to directory: %v", err)
	}

	for {
		dbPath := dir + string(os.PathSeparator) + dbFileName
		// fmt.Println("Checking", dbPath)
		if _, err := os.Stat(dbPath); err == nil {
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				return nil, fmt.Errorf("unable to open %s. Is the file corrupted? %v", dbFileName, err)
			}
			return db, nil
		}

		parentDir := dir + string(os.PathSeparator) + ".."
		parentDir, err = filepath.Abs(parentDir)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve parent directory: %v", err)
		}
		if parentDir == dir {
			return nil, fmt.Errorf("%s file not found\n%s to create a new mindtick", dbFileName, messages.ColorizeStr("mindtick new", messages.BrightGreen))
		}
		dir = parentDir
	}
	// return nil, fmt.Errorf("mindtick file not found. you shouldn't have gotten here")
}

// `mindtick new` command
func New() error {
	// Verify if file exists
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		// do nothing
	} else {
		return fmt.Errorf("%s already exists", dbFileName)
	}

	// check if .gitignore exists and append store.mindtick
	if _, err := os.Stat(".gitignore"); err == nil {
		gitignore, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .gitignore: %v", err)
		}
		defer gitignore.Close()
		if _, err := gitignore.WriteString("\n" + dbFileName + "\n"); err != nil {
			return fmt.Errorf("failed to write to .gitignore: %v", err)
		}
	}

	file, err := os.Create(dbFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s %v", dbFileName, err)
	}
	defer file.Close()

	// setup database schema
	db, err := LoadMindtick()
	if err != nil {
		return err
	}
	createSchema(db)

	return fmt.Errorf(messages.ColorizeStr("mindtick intialized", messages.BrightPurple))
}

func Delete() error {
	err := os.Remove(dbFileName)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %v", dbFileName, err)
	}
	return fmt.Errorf(messages.ColorizeStr("mindtick deleted", messages.BrightPurple))
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
	RangeToTime = map[Range]time.Time{
		TODAY:     time.Now().Truncate(24 * time.Hour),
		YESTERDAY: time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour),
		WEEK:      time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour),
		MONTH:     time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour),
	}
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
		SQLstmt = "SELECT * FROM messages WHERE timestamp >= ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, RangeToTime[rangeType])
	}

	if rangeType == ANYTIME && tag != messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages WHERE msgtype = ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, tag)
	}

	if rangeType != ANYTIME && tag != messages.ANYTAG {
		SQLstmt = "SELECT * FROM messages WHERE msgtype = ? AND timestamp >= ? ORDER BY timestamp"
		rows, err = db.Query(SQLstmt, tag, RangeToTime[rangeType])
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
