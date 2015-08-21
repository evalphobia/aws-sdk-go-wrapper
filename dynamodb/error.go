package dynamodb

import (
	"strings"
)

const errPrefix = "[DynamoDB] "

type DynamoError struct {
	errList []string
}

func NewError(msg string) *DynamoError {
	return &DynamoError{
		errList: []string{msg},
	}
}

func NewErrorList(msgList []string) *DynamoError {
	return &DynamoError{
		errList: msgList,
	}
}

func (e DynamoError) Error() string {
	return errPrefix + strings.Join(e.errList, " || ")
}

func (e *DynamoError) AddMessage(msg string) {
	e.errList = append(e.errList, msg)
}

func (e *DynamoError) HasError() bool {
	return len(e.errList) != 0
}
