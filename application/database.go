package application

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Создает базы данных по заданной конфигурации, если базы данных
// уже существовали, к ним произойдет подключение
func NewDatabases(config *ApplicationConfig) []*sql.DB {

	genericPath := config.Storage.Path
	volumes := config.Storage.Volumes

	databases := make([]*sql.DB, 0)

	if DirectoryExists(genericPath) {
		var i int32 = 1
		for i <= volumes {
			currentDbUrl := fmt.Sprintf("file:%sstore%d.db?_journal_mode=WAL", genericPath, i)
			db, err := sql.Open("sqlite3", currentDbUrl)
			if err != nil {
				panic(err)
			}
			initializeDb(db)
			databases = append(databases, db)
			i += 1
		}
	} else {
		log.Fatalf("Storage directory does not exist, check setting `storage-path`: %v", genericPath)
	}

	return databases
}

func initializeDb(database *sql.DB) {
	initializeFilesTable(database)
	initializeMetaInfoTable(database)
}

func initializeFilesTable(database *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS files (
		id int64 PRIMARY KEY AUTOINCREMENT, 
		mime_type text NOT NULL, 
		content blob, 
		file_size int64 NOT NULL,
		status integer NOT NULL);`
	_, err := database.Exec(query)
	if err != nil {
		panic(err)
	}
}

func initializeMetaInfoTable(database *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS meta_info (
		total_size int64 NOT NULL DEFAULT 0,
		compacting integer NOT NULL DEFAULT 0);`
	_, err := database.Exec(query)
	if err != nil {
		panic(err)
	}
}
