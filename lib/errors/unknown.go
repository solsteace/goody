package errors

import "fmt"

type Unknown struct {
	err error
}

func (e Unknown) Error() string {
	return fmt.Sprintf("Unknown error: %v", e.err)
}
