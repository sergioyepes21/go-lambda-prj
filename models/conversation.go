package models

import (
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func conversationTE() TableEntity {
	return TableEntity{
		PK: "$ConversationId",
		SK: "#Aametadata",
	}
}

type Conversation struct {
	TE             TableEntity
	ConversationId string `validator:"required" json:"conversationId"`
	Name           string `validator:"required" json:"name"`
	CreatedAt      int64  `validator:"required,number" json:"createdAt"`
}

func NewConversation(conversationId string, name string, createdAt int64) *Conversation {
	return &Conversation{
		TE:             conversationTE(),
		ConversationId: conversationId,
		CreatedAt:      createdAt,
		Name:           name,
	}
}

func NewConversationById(conversationId string) *Conversation {
	return &Conversation{
		TE:             conversationTE(),
		ConversationId: conversationId,
	}
}

func GetConversation(conversationId string, repository Repository) (Conversation, error) {
	entity := NewConversationById(conversationId)
	modelReflectValue := reflect.ValueOf(entity)

	pk, err := TransformPk(modelReflectValue, entity.TE)
	if err != nil {
		return Conversation{}, err
	}

	resultItems, err := repository.GetItemsByPk(pk, []string{"ConversationId", "Name", "CreatedAt"})
	if err != nil {
		return Conversation{}, err
	}

	conversationList := fromDynamoResultTo(resultItems)
	if len(conversationList) == 0 {
		return Conversation{}, nil
	}

	return conversationList[0], nil
}

func fromDynamoResultTo(attributeValueList []map[string]*dynamodb.AttributeValue) []Conversation {
	var conversationList []Conversation
	for i := 0; i < len(attributeValueList); i++ {
		attributeValue := attributeValueList[i]

		createdAt, _ := strconv.ParseInt(*attributeValue["CreatedAt"].N, 10, 64)
		conversation := NewConversation(*attributeValue["ConversationId"].S, *attributeValue["Name"].S, createdAt)

		conversationList = append(conversationList, *conversation)
	}

	return conversationList
}
