package utils

import "fmt"

type ErrWrongUrlValues struct {
	valueName string
}

func (err *ErrWrongUrlValues) Error() string {
	return fmt.Sprintf("wrong url value=%s", err.valueName)
}
