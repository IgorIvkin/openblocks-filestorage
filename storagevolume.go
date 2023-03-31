package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type StorageVolume struct {
	Index      int
	Database   *sql.DB
	TotalSize  int64
	Compacting bool
	Busy       bool
	sync.Mutex
}

type MetaInfo struct {
	TotalSize  int64
	Compacting int32
}

// Инициализирует доступные тома
func NewStorageVolumes(databases []*sql.DB) []*StorageVolume {
	storageVolumes := make([]*StorageVolume, 0)

	for index, database := range databases {
		volume := initializeStorageVolume(database, index)
		storageVolumes = append(storageVolumes, volume)
	}

	return storageVolumes
}

// Сохраняет файл в выбранном томе
func (volume *StorageVolume) StoreFile(content []byte) (string, error) {
	volume.Lock()
	defer volume.Unlock()

	volume.Busy = true

	lastInsertId, err := storeContentOfFile(volume.Database, content)
	if err != nil {
		volume.Busy = false
		return "", err
	}

	volume.Busy = false
	return fmt.Sprintf("%d-%d", volume.Index, lastInsertId), nil
}

func initializeStorageVolume(database *sql.DB, index int) *StorageVolume {
	rows, err := database.Query("SELECT total_size, compacting FROM meta_info")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var found bool = false
	var metaInfo MetaInfo
	for rows.Next() {
		metaInfo = MetaInfo{}
		err := rows.Scan(&metaInfo.TotalSize, &metaInfo.Compacting)
		if err != nil {
			panic(err)
		}
		found = true
	}

	var storageVolume StorageVolume
	if found {
		storageVolume = StorageVolume{
			Index:      index,
			Database:   database,
			TotalSize:  metaInfo.TotalSize,
			Compacting: intToBool(metaInfo.Compacting),
			Busy:       false,
		}
	} else {
		insertMetaInfo(database)
		storageVolume = StorageVolume{
			Index:      index,
			Database:   database,
			TotalSize:  0,
			Compacting: false,
			Busy:       false,
		}
	}
	return &storageVolume
}

func storeContentOfFile(database *sql.DB, content []byte) (int64, error) {
	transaction, _ := database.Begin()

	// Вставляем файл в хранилище
	statement, err := transaction.Prepare(`INSERT INTO files(mime_type, content, file_size, status) VALUES (?, ?, ?, ?)`)
	if err != nil {
		transaction.Rollback()
		log.Printf("Cannot prepare query to insert file, reason: %v", err)
		return 0, errors.New("cannot insert file")
	}
	result, err := statement.Exec("", content, len(content), 1)
	if err != nil {
		transaction.Rollback()
		log.Printf("Cannot insert file, reason: %v", err)
		return 0, errors.New("cannot insert file")
	}

	// Обновляем мета-информацию в томе
	updateMetaInfo(transaction, len(content))

	// Получаем последний сгенерированный идентификатор
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		log.Printf("Cannot update meta-info, reason: %v", err)
		return 0, errors.New("cannot update meta-info")
	}
	transaction.Commit()
	return lastInsertId, nil
}

func updateMetaInfo(transaction *sql.Tx, contentSize int) {
	statement, err := transaction.Prepare(`UPDATE meta_info SET total_size = total_size + ?`)
	if err != nil {
		transaction.Rollback()
		log.Printf("Cannot prepare query to update meta-info, reason: %v", err)
		return
	}
	_, err2 := statement.Exec(contentSize)
	if err2 != nil {
		transaction.Rollback()
	}
}

func insertMetaInfo(database *sql.DB) {
	query := `INSERT INTO meta_info(total_size, compacting) VALUES ($1, $2)`
	_, err := database.Exec(query, 0, 0)
	if err != nil {
		panic(err)
	}
}

func intToBool(value int32) bool {
	return value == 1
}
