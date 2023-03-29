package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	application := NewApplication()
	defer application.CloseDb()

	fmt.Println("Test")

	router := gin.Default()
	router.MaxMultipartMemory = 10 << 20 // ограничиваем максимальный размер файла в 10 мегабайтов

	router.POST("/api/v1/store", func(context *gin.Context) {
		ProcessStoreFile(application, context)
	})

	router.Run(":8903")
}
