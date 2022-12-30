package dataservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	TAG_NAME      = "dyorm"
	TAG_SEPARATOR = "$"
)

type Repository struct {
	client dynamodbiface.DynamoDBAPI
}

func NewRepository(client dynamodbiface.DynamoDBAPI) *Repository {
	return &Repository{client: client}
}

// func (r *Repository) transformPk(tableEntity *models.TableEntity) (string, error) {
// 	return r.transform(tableEntity, "PK", false)
// }

func (r *Repository) GetItemsByPk(pk string, projectionExp []string) ([]map[string]*dynamodb.AttributeValue, error) {
	// newKey, err := r.transformPk(tableEntity)

	// if err != nil {
	// 	return nil, err
	// }
	keyCond := expression.Key("PK").Equal(expression.Value(pk))
	proj := expression.NamesList(expression.Name("aName"), expression.Name("anotherName"), expression.Name("oneOtherName"))
	builder := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj)
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	queryInput := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String("conversation-table"),
	}

	svc := dynamodb.New(session.New())
	result, err := svc.Query(&queryInput)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
