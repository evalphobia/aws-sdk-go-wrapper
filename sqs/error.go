// SQS Queue

package sqs

import (
	"strings"
)

const errPrefix = "[SQS] "

type SQSError struct {
	errList []string
}

func NewError(msg string) *SQSError {
	return &SQSError{
		errList: []string{msg},
	}
}

func NewErrorList(msgList []string) *SQSError {
	return &SQSError{
		errList: msgList,
	}
}

func (e SQSError) Error() string {
	return errPrefix + strings.Join(e.errList, " || ")
}

func (e *SQSError) Add(err error) {
	e.errList = append(e.errList, err.Error())
}

func (e *SQSError) AddMessage(msg string) {
	e.errList = append(e.errList, msg)
}

func (e *SQSError) HasError() bool {
	return len(e.errList) != 0
}
