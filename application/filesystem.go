package application

import (
	"log"
	"os"
)

// Проверяет существование директории. Если директории не существует,
// возвращается false, при других ошибках доступа к директории, происходит
// падение приложения
func DirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}
