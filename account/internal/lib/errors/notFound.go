package errors

type NotFound struct{}

func (e NotFound) Error() string {
	return "Record not found"
}
