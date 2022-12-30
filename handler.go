package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	dataservice "github.com/sergioyepes21/go-lambda-prj/data-service"
	"github.com/sergioyepes21/go-lambda-prj/models"
)

type MyEvent struct {
	Name string `json:"name"`
}

type MyResponse struct {
	Message string `json:"Message"`
}

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dataService := dataservice.NewDataService()
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }
	conversationM := models.Conversation{
		ConversationId: "conv123",
	}
	projectionExp := []string{"ConversationId", "Name"}
	result, err := dataService.Repository.GetItemsByPk(&conversationM.TE, projectionExp)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	conversationList := conversationM.FromDynamoResultTo(result)
	conversationListJson, err := json.Marshal(conversationList)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body: string(conversationListJson),
	}, nil
}

func main() {
	lambda.Start(handler)
}
