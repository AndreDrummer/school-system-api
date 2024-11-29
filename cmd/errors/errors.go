package apperrors

import "fmt"

type NotFoundError struct{}

type JsonEncodingError struct {
	Type any
	Err  error
}

type JsonDecodingError struct {
	Type any
	Err  error
}

func (n *NotFoundError) Error() string {
	return "Document not found"
}

func (j *JsonEncodingError) Error() string {
	return fmt.Sprintf("Error converting data of type %v to JSON, %s", j.Type, j.Err.Error())
}

func (j *JsonDecodingError) Error() string {
	return fmt.Sprintf("Error converting JSON to %v: %s", j.Type, j.Err.Error())
}
