package errors

import (
	"fmt"
	"strings"
)

const defaultSeparator = " || "

// Errors is struct for storing multiple errors.
type Errors struct {
	serviceName string
	separator   string
	list        []string
}

// NewErrors returns initialized *Errors.
func NewErrors(name string) *Errors {
	return &Errors{
		serviceName: name,
		separator:   defaultSeparator,
	}
}

func (e Errors) Error() string {
	return fmt.Sprintf("[%s] %s", e.serviceName, strings.Join(e.list, e.separator))
}

// Add adds error into the error list.
func (e *Errors) Add(err error) {
	e.list = append(e.list, err.Error())
}

// AddMessage adds error message into the error list.
func (e *Errors) AddMessage(msg string) {
	e.list = append(e.list, msg)
}

// HasError checks the error list contains one or more error or not.
func (e *Errors) HasError() bool {
	return len(e.list) != 0
}
