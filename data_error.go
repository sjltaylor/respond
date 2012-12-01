package respond

import (
	"bytes"
	"fmt"
	"strings"
)

type DataError struct {
	Message string
	Errors  map[string][]string
}

func NewDataError() *DataError {

	return &DataError{
		Errors: make(map[string][]string),
	}
}

func NewDataErrorWithMessage(message string) *DataError {

	e := NewDataError()

	e.Message = message

	return e
}

func (de *DataError) Add(key, errorMessage string) {

	errorMessages, ok := de.Errors[key]

	if !ok {
		errorMessages = []string{}
	}

	errorMessages = append(errorMessages, errorMessage)

	de.Errors[key] = errorMessages
}

func (de *DataError) Error() string {

	buffer := bytes.NewBufferString(fmt.Sprintf("%s:\n", de.Message))

	for k, _ := range de.Errors {
		buffer.WriteString(fmt.Sprintf("\t%s: %s\n", k, de.ErrorsFor(k)))
	}

	return buffer.String()
}

func (de *DataError) ErrorsFor(key string) string {

	errorMessages, ok := de.Errors[key]

	if !ok {
		errorMessages = []string{}
	}

	return strings.Join(errorMessages, ", ")
}

func (de *DataError) HasDetailsFor(key string) bool {

	_, hasDetailsForKey := de.Errors[key]

	return !hasDetailsForKey
}

func (de *DataError) HasDetails() bool {
	for _, errs := range de.Errors {

		if len(errs) > 0 {
			return true
		}
	}
	return false
}

func (de *DataError) HTTPStatusCode() int {
	return 422
}
