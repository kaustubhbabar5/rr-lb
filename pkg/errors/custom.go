package errors

import "fmt"

type NotFound struct {
	Kind   string
	Object string
	Msg    string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("entity: %s, message: %s", e.Object, e.Msg)
}

func NewNotFound(object, msg string) *NotFound {
	return &NotFound{
		Kind:   "not_found",
		Object: object,
		Msg:    msg,
	}
}
