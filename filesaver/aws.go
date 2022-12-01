package filesaver

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/VladimirZaets/distribution-portal/metadata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "distribution-portal"
)

type AwsStorage struct {
	storageType string
	c           *gin.Context
}

func NewAWSFileSaver(c *gin.Context) *AwsStorage {
	return &AwsStorage{
		storageType: "aws",
		c:           c,
	}
}

func (ls *AwsStorage) GetType() string {
	return ls.storageType
}

func (ls *AwsStorage) Save(file *multipart.FileHeader) error {
	ttype := true
	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION), CredentialsChainVerboseErrors: &ttype})

	if err != nil {
		log.Fatal(err)
		return err
	}
	err = uploadFile(session, file)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (ls *AwsStorage) Get(m *metadata.Metadata) (*File, error) {
	fmt.Println("INSSSS")
	ttype := true
	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION), CredentialsChainVerboseErrors: &ttype})
	if err != nil {
		return nil, err
	}
	res, err := s3.New(session).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(m.Name),
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("res.ContentLength1", res.ContentLength)
	fmt.Println("res.ContentLength2", res.ContentType)
	return &File{
		ContentLength: *res.ContentLength,
		ContentType:   *res.ContentType,
		Reader:        res.Body,
	}, nil
}

func uploadFile(session *session.Session, upFile *multipart.FileHeader) error {

	var fileSize int64 = upFile.Size
	fileBuffer := make([]byte, fileSize)
	fileContent, _ := upFile.Open()
	fileContent.Read(fileBuffer)

	id, err := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
		Key:                  aws.String(upFile.Filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	fmt.Println(id)
	return err
}
