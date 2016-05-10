package trace

import (
	"fmt"
)

type Error struct {
	BaseError interface{}
	Stack     []byte
}

func (a *Error) Error() string {
	return fmt.Sprintf("%s at:\n%s", a.BaseError, string(a.Stack))
}
