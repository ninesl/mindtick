package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

	return nil
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

func Delete() error {
	err := os.Remove(dbFileName)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}
	return nil
}

func AddMessage(db *sql.DB, message messages.Message) error {
	_, err := db.Exec("INSERT INTO messages (timestamp, msg, msgtype) VALUES (?, ?, ?)", message.Timestamp, message.Msg, message.MsgType)
	if err != nil {
		return fmt.Errorf("unable to add message: %v", err)
	}
	return nil
}

func GetMessages(db *sql.DB) ([]messages.Message, error) {
	rows, err := db.Query("SELECT * FROM messages ORDER BY timestamp DESC")
	if err != nil {
		return nil, fmt.Errorf("unable to select messages: %v", err)
	}
	defer rows.Close()

	var msgs []messages.Message
	for rows.Next() {
		var msg messages.Message
		err := rows.Scan(&msg.ID, &msg.Timestamp, &msg.Msg, &msg.MsgType)
		if err != nil {
			return nil, fmt.Errorf("unable to read message: %v", err)
		}
		msgs = append(msgs, msg)
	}
	if len(msgs) == 0 {
		return nil, fmt.Errorf("no messages found")
	}
	return msgs, nil
}
