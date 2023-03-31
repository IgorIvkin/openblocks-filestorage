package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Application struct {
	Config    *ApplicationConfig
	Databases []*sql.DB
	Volumes   []*StorageVolume
}

type ApplicationConfig struct {
	Storage ApplicationConfigStorage `yaml:"storage"`
}

type ApplicationConfigStorage struct {
	Path          string `yaml:"path"`
	Volumes       int32  `yaml:"volumes"`
	MaxVolumeSize int32  `yaml:"max-volume-size"`
}

// Возвращает экземпляр приложения, в нём представлены конфигурация, список баз данных
// и список томов для хранения файлов
func NewApplication() *Application {
	rand.Seed(time.Now().Unix())
	config := getConfig()
	databases := NewDatabases(config)
	app := Application{
		Config:    config,
		Databases: databases,
		Volumes:   NewStorageVolumes(databases),
	}
	return &app
}

// Закрывает соединения с базаыми данных
func (app *Application) CloseDb() {
	for _, database := range app.Databases {
		database.Close()
	}
}

// Выбирает подходящий том, который не занят обработкой данных
func (app *Application) ChooseIdleVolume() *StorageVolume {
	var selectedVolume *StorageVolume

	// Выбираем один случайный том, если он оказался занят работой,
	// прокручиваем все тома, пока не найдем свободный
	volumesCount := len(app.Volumes)
	if volumesCount == 0 {
		log.Fatal("Cannot choose idle volume, no volumes presented, check that you specified correct \"storage.path\" parameter")
	}
	selectedVolume = app.Volumes[rand.Intn(volumesCount)]
	if selectedVolume.Busy || selectedVolume.Compacting {
		for {
			var found bool
			for _, volume := range app.Volumes {
				if !volume.Busy && !volume.Compacting {
					found = true
					selectedVolume = volume
					break
				}
			}
			if found {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	return selectedVolume
}

// Возвращает конфигурацию, заданную в yml-файле приложения,
// в ней можно задать множественные рейт-лимитеры.
func getConfig() *ApplicationConfig {
	configFile, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var config ApplicationConfig
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
