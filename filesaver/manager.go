package filesaver

import (
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

type filesaver interface {
	Save(file *multipart.FileHeader) error
	GetType() string
}

type Manager struct {
	C *gin.Context
}

func (fs *Manager) Get() filesaver {
	switch os.Getenv("DIST_PORTAL_ENV") {
	case "prod":
		return NewAWSFileSaver(fs.C)
	case "dev":
		return NewLocalFileSaver(fs.C)
	}

	return NewLocalFileSaver(fs.C)
}
