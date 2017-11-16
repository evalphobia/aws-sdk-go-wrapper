package rekognition

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErrorInvalidImageEncoding(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expected bool
		text     string
	}{
		{true, "InvalidImageFormatException: Invalid image encoding"},
		{true, "status code: 400, request id: 00ff00ff-00ff-00ff-00ff-00ff00ff00ff, InvalidImageFormatException: Invalid image encoding"},
		{true, "InvalidImageFormatException: Invalid image encoding, error"},
		{true, "error, InvalidImageFormatException: Invalid image encoding, error"},
		{false, "No errors"},
		{false, "InvalidImageFormatException: Invalid image encodin"},
		{false, "nvalidImageFormatException: Invalid image encoding"},
		{false, ""},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		err := errors.New(tt.text)
		a.Equal(tt.expected, IsErrorInvalidImageEncoding(err), target)
	}
	a.Equal(false, IsErrorInvalidImageEncoding(nil), "When error=nil")
}
