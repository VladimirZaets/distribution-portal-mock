package filesaver

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gin-gonic/gin"
)

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "mypersonaltestbucket"
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
	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = uploadFile(session, file)
	file.Open()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func uploadFile(session *session.Session, upFile *multipart.FileHeader) error {

	var fileSize int64 = upFile.Size
	fileBuffer := make([]byte, fileSize)
	fileContent, _ := upFile.Open()
	fileContent.Read(fileBuffer)

	_, err := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
		Key:                  aws.String(upFile.Filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
