package metadata

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Manager struct {
	C *gin.Context
}

type MetadataList map[string]*Metadata

type MetadataInterface interface {
	GetList() (MetadataList, error)
	Get(*Metadata) (*Metadata, error)
	Set(*Metadata) error
	Delete(*Metadata) error
	Update(*Metadata) error
}

type Metadata struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func (m *Manager) Get() MetadataInterface {
	switch os.Getenv("DIST_PORTAL_ENV") {
	case "prod":
		return NewAWSMetadata()
	case "dev":
		return NewLocalMetadata()
	}

	return NewLocalMetadata()
}
