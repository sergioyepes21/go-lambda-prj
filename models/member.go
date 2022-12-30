package models

import (
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func memberTE() TableEntity {
	return TableEntity{
		PK: "$UserId",
		SK: "$AddedAt-$ConversationId",
	}
}

type Member struct {
	TE        TableEntity
	MemberId  string `validator:"required" json:"memberId"`
	UserId    string `validator:"required" json:"conversationId"`
	AddedAt   int64  `validator:"required,number" json:"addedAt"`
	CreatedAt int64  `validator:"required,number" json:"createdAt"`
}

func NewMember(memberId string, userId string, addedAt int64, createdAt int64) *Member {
	return &Member{
		TE:        memberTE(),
		MemberId:  memberId,
		CreatedAt: createdAt,
		AddedAt:   addedAt,
	}
}

func NewMemberByUserId(userId string) *Member {
	return &Member{
		TE:     memberTE(),
		UserId: userId,
	}
}

func GetUserMembers(userId string, repository Repository) ([]Member, error) {
	entity := NewMemberByUserId(userId)
	modelReflectValue := reflect.ValueOf(entity)
	pk, err := TransformPk(modelReflectValue, entity.TE)
	if err != nil {
		return []Member{}, err
	}
	resultItems, err := repository.GetItemsByPk(pk, []string{"UserId", "MemberId", "CreatedAt", "AddedAt"})
	if err != nil {
		return []Member{}, err
	}

	memberList := fromDynamoResultTo(resultItems)

	return memberList, nil
}

func fromDynamoResultTo(attributeValueList []map[string]*dynamodb.AttributeValue) []Member {
	var entityList []Member
	for i := 0; i < len(attributeValueList); i++ {
		attributeValue := attributeValueList[i]
		createdAt, _ := strconv.ParseInt(*attributeValue["CreatedAt"].N, 10, 64)
		addedAt, _ := strconv.ParseInt(*attributeValue["AddedAt"].N, 10, 64)
		entity := NewMember(*attributeValue["MemberId"].S, *attributeValue["UserId"].S, addedAt, createdAt)
		entityList = append(entityList, *entity)
	}
	return entityList
}
