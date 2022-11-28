package filesaver

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

type LocalStorage struct {
	storageType string
	c           *gin.Context
}

func NewLocalFileSaver(c *gin.Context) *LocalStorage {
	return &LocalStorage{
		storageType: "local",
		c:           c,
	}
}

func (ls *LocalStorage) GetType() string {
	return ls.storageType
}

func (ls *LocalStorage) Save(file *multipart.FileHeader) error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	ls.c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", path, file.Filename))
	return nil
}
