package dataservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sergioyepes21/go-lambda-prj/models"
)

type IRepository interface {
	GetItemsByPk(models.TableEntity, []string) ([]models.TableEntity, error)
}

type DataService struct {
	Repository Repository
}

func NewDataService() *DataService {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	client := dynamodb.New(sess)
	return &DataService{
		Repository: *NewRepository(client),
	}
}
