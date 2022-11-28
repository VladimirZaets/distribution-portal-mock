package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VladimirZaets/distribution-portal/filesaver"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		fsaver := &filesaver.Manager{
			C: c,
		}
		saverSdk := fsaver.Get()
		fmt.Println(saverSdk.GetType())
		log.Println(file.Filename)
		saverSdk.Save(file)
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	r.Run()
}
