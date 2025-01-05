package store

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
	//https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc
)

const dbFileName = "store.mindtick"

func loadMindtick() *sql.DB {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		dbPath := dir + string(os.PathSeparator) + dbFileName
		if _, err := os.Stat(dbPath); err == nil {
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				log.Fatalf("found mindtick, unable to load. Is the file corrupted? - %v", err)
			}
			return db
		}

		parentDir := dir + string(os.PathSeparator) + ".."
		if parentDir == dir {
			log.Fatalf("mindtick file not found. Make sure to run `mindtick new` first")
		}
		dir = parentDir
	}
}

// `mindtick new` command
func New() {
	// Verify if file exists
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		log.Printf("Creating new mindtick in currrent directory\n")
		// fmt.Println()
	} else {
		log.Printf("%s already exists in this directory", dbFileName)
		return
	}

	// else create new file
	file, err := os.Create(dbFileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()
	// setup database schema
	db := loadMindtick()
	createSchema(db)

	// log to user that file is created
	log.Printf("New mindtick created under %s\n", dbFileName)
}

func createSchema(db *sql.DB) {
	// schema is
	// messages
	// 	id - serial
	// 	timestamp - time
	// 	msg - string
	// 	msgtype - int enums

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		msg TEXT,
		msgtype INT
	)`)
	if err != nil {
		log.Fatalf("Failed to create mindtick schema: %v", err)
	}

}
