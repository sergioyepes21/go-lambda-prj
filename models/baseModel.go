package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	TAG_SEPARATOR = "$"
)

type Repository interface {
	GetItemsByPk(pk string, projectionExp []string) ([]map[string]*dynamodb.AttributeValue, error)
}

func TransformPk(modelReflectValue reflect.Value, tableEntity TableEntity) (string, error) {
	return TransformKey(modelReflectValue, tableEntity, "PK", true)
}

func TransformKey(modelReflectValue reflect.Value, tableEntity TableEntity, keyName string, shouldHaveAllAtts bool) (string, error) {
	valueOfTE := reflect.ValueOf(tableEntity)
	keyValue := valueOfTE.FieldByName(keyName).String()
	keyArr := strings.Split(keyValue, TAG_SEPARATOR)
	var newVariableArr []string
	// newSValue := reflect.ValueOf(s)
	for i := 0; i < len(keyArr); i++ {
		att := keyArr[i]
		if att[0] == '$' {
			attValue := modelReflectValue.FieldByName(att)
			if shouldHaveAllAtts && attValue.IsZero() {
				return "", fmt.Errorf("the model has missing required parameters: %s", att)
			}
			newVariableArr = append(newVariableArr, attValue.String())
		} else {
			newVariableArr = append(newVariableArr, att)
		}
	}
	newVariable := strings.Join(newVariableArr, TAG_SEPARATOR)
	return newVariable, nil
}

type TableEntity struct {
	PK string
	SK string
}
