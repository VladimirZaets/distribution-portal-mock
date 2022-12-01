package metadata

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	aws_region         string = "us-east-1"
	aws_dynamodb_table string = "distribution-portal"
)

type AwsMetadata struct {
	metadataType string
	awsSession   *session.Session
	dynamodbSdk  *dynamodb.DynamoDB
}

func NewAWSMetadata() *AwsMetadata {
	credentialsChainVerboseErrors := true
	session, err := session.NewSession(&aws.Config{Region: aws.String(aws_region), CredentialsChainVerboseErrors: &credentialsChainVerboseErrors})
	if err != nil {
		log.Fatal(err)
	}

	return &AwsMetadata{
		metadataType: "local",
		awsSession:   session,
		dynamodbSdk:  dynamodb.New(session),
	}
}

func (am *AwsMetadata) Set(m *Metadata) error {
	item, err := dynamodbattribute.MarshalMap(m)
	if err != nil {
		return err
	}
	_, err = am.dynamodbSdk.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(aws_dynamodb_table),
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (am *AwsMetadata) Get(m *Metadata) (*Metadata, error) {
	res, err := am.dynamodbSdk.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(aws_dynamodb_table),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(m.Name),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	data := &Metadata{}
	err = dynamodbattribute.UnmarshalMap(res.Item, data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return data, nil
}

func (am *AwsMetadata) GetList() (MetadataList, error) {
	res, err := am.dynamodbSdk.Scan(&dynamodb.ScanInput{
		TableName: aws.String(aws_dynamodb_table),
	})
	if err != nil {
		return nil, err
	}

	data := map[string]*Metadata{}
	for _, i := range res.Items {
		item := &Metadata{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		data[item.Name] = item
	}

	return data, nil
}

func (am *AwsMetadata) Update(m *Metadata) error {
	return fmt.Errorf("Error")
}

func (am *AwsMetadata) Delete(m *Metadata) error {
	return fmt.Errorf("Error")
}

func (am *AwsMetadata) GetType() string {
	return am.metadataType
}
