package main

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/VladimirZaets/distribution-portal/filesaver"
	"github.com/VladimirZaets/distribution-portal/metadata"
	"github.com/gin-gonic/gin"
)

type UploadRequestData struct {
	Name    string                `form:"name" binding:"required"`
	Version string                `form:"version" binding:"required"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}

func main() {
	r := gin.Default()
	saver := (&filesaver.Manager{}).Get()
	meta := (&metadata.Manager{}).Get()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong" + os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"),
		})
	})
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.GET("/packages", func(c *gin.Context) {
		data, err := meta.GetList()
		if err != nil {
			fmt.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
			})
		}
		fmt.Println(data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
			})
		}
		fmt.Println(data)

		c.JSON(http.StatusOK, gin.H{
			"body": string(jsonData),
		})
	})

	r.GET("/package/:name", func(c *gin.Context) {
		name := c.Param("name")
		// data, err := meta.Get(&metadata.Metadata{
		// 	Name: name,
		// })
		file, err := saver.Get(&metadata.Metadata{
			Name: name,
		})
		if err != nil {
			fmt.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
			})
		}
		fmt.Println(saver.GetType())
		fmt.Println(file)
		fmt.Println("file.ContentLength", file.ContentLength, "file.ContentType", file.ContentType)
		c.DataFromReader(200, file.ContentLength, file.ContentType, file.Reader, nil)

		// data, err := meta.Get(&metadata.Metadata{
		// 	Name: name,
		// })
		// if err != nil {
		// 	fmt.Errorf(err.Error())
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": true,
		// 	})
		// }

		// jsonData, err := json.Marshal(data)
		// if err != nil {
		// 	fmt.Errorf(err.Error())
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": true,
		// 	})
		// }
		// c.JSON(http.StatusOK, gin.H{
		// 	"body": string(jsonData),
		// })
	})
	r.POST("/package", func(c *gin.Context) {
		data := &UploadRequestData{}
		err := c.Bind(data)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		meta.Set(&metadata.Metadata{
			Name:    data.Name,
			Version: data.Version,
		})
		saver.Save(data.File)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", data.File.Filename))
	})
	r.Run()
}
