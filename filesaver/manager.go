package filesaver

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/VladimirZaets/distribution-portal/metadata"
	"github.com/gin-gonic/gin"
)

type filesaver interface {
	Save(file *multipart.FileHeader) error
	Get(m *metadata.Metadata) (*File, error)
	GetType() string
}

type File struct {
	ContentType   string
	ContentLength int64
	Reader        io.Reader
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
