package kinesis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		iterator string
		expected string
	}{
		{"A", "A"},
		{"ABC", "ABC"},
		{"", "LATEST"},
		{"LATEST", "LATEST"},
		{string(IteratorTypeLatest), "LATEST"},
		{string(IteratorTypeTrimHorizon), "TRIM_HORIZON"},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		it := IteratorType(tt.iterator)
		assert.Equal(tt.iterator, string(it), target)
		assert.Equal(tt.expected, it.String(), target)
	}
}

func TestIsEmpty(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		iterator string
		expected bool
	}{
		{"", true},
		{"A", false},
		{"ABC", false},
		{"LATEST", false},
		{string(IteratorTypeLatest), false},
		{string(IteratorTypeTrimHorizon), false},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		it := IteratorType(tt.iterator)
		assert.Equal(tt.expected, it.isEmpty(), target)
	}
}
