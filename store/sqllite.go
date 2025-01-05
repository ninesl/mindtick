package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
	//https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc
)

const dbFileName = "store.mindtick"

func loadMindtick() (*sql.DB, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get access to directory: %v", err)
	}

	for {
		dbPath := dir + string(os.PathSeparator) + dbFileName
		if _, err := os.Stat(dbPath); err == nil {
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				return nil, fmt.Errorf("unable to open %s. Is the file corrupted? %v", dbFileName, err)
			}
			return db, nil
		}

		parentDir := dir + string(os.PathSeparator) + ".."
		if parentDir == dir {
			return nil, fmt.Errorf("%s file not found", dbFileName)
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

	file, err := os.Create(dbFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s %v", dbFileName, err)
	}
	defer file.Close()

	// setup database schema
	db, err := loadMindtick()
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
