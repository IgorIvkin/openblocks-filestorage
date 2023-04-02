package application

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Application struct {
	Config    *ApplicationConfig
	Databases []*sql.DB
	Volumes   []*StorageVolume
}

type ApplicationConfig struct {
	Storage ApplicationConfigStorage   `yaml:"storage"`
	General ApplicationGeneralSettings `yaml:"general"`
}

type ApplicationConfigStorage struct {
	Path          string `yaml:"path"`
	Volumes       int32  `yaml:"volumes"`
	MaxVolumeSize int32  `yaml:"max-volume-size"`
}

type ApplicationGeneralSettings struct {
	UseJwtAuth bool     `yaml:"use-jwt-auth"`
	JwksUrl    string   `yaml:"jwks-url"`
	ExceptUrls []string `yaml:"except-urls"`
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

// Получает файл из тома, если такой файл не найден, возвращается ошибка.
// Если файл найден, функция возвращает его содержимое в виде массива байтов,
// а также MIME-тип хранимых данных, что позволяет сформировать корректный ответ
func (app *Application) GetFile(fileId string) ([]byte, string, error) {
	parts := strings.Split(fileId, "-")
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("incorrect file ID %s", fileId)
	}

	volumeId, err := strToInt32(parts[0])
	if err != nil {
		return nil, "", err
	}

	entityId, err := strToInt64(parts[1])
	if err != nil {
		return nil, "", err
	}

	for index, volume := range app.Volumes {
		if index == volumeId {
			return volume.GetFile(entityId)
		}
	}

	return nil, "", fmt.Errorf("file not found by ID %s", fileId)
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

func strToInt32(value string) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("string value is incorrect integer: %s", value)
	}
	return number, nil
}

func strToInt64(value string) (int64, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("string value is incorrect int64: %s", value)
	}
	return number, nil
}
