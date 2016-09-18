package sqs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)

	err := NewError("error message")
	assert.Len(err.errList, 1)
	assert.Equal("error message", err.errList[0])
}

func TestNewErrorList(t *testing.T) {
	assert := assert.New(t)

	errList := []string{"e1", "e2", "e3"}

	err := NewErrorList(errList)
	assert.Len(err.errList, 3)
	assert.Equal("e1", err.errList[0])
	assert.Equal("e2", err.errList[1])
	assert.Equal("e3", err.errList[2])
}

func TestError(t *testing.T) {
	assert := assert.New(t)

	list := []string{"e1", "e2", "e3"}
	err := NewErrorList(list)
	assert.Equal("[SQS] e1 || e2 || e3", err.Error())
}

func TestErrorAdd(t *testing.T) {
	assert := assert.New(t)

	err := NewError("e0")
	assert.Len(err.errList, 1)

	err.Add(errors.New("e1"))
	assert.Len(err.errList, 2)
	assert.Equal("[SQS] e0 || e1", err.Error())
}

func TestErrorAddMessage(t *testing.T) {
	assert := assert.New(t)

	err := NewError("e0")
	assert.Len(err.errList, 1)

	err.AddMessage("e1")
	assert.Len(err.errList, 2)
	assert.Equal("[SQS] e0 || e1", err.Error())
}

func TestErrorHasError(t *testing.T) {
	assert := assert.New(t)

	err := NewError("e0")
	assert.True(err.HasError())

	err = &SQSError{}
	assert.False(err.HasError())
}
