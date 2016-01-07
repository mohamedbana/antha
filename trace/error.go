package trace

import (
	"fmt"
)

type goError struct {
	BaseError interface{}
	Stack     []byte
}

func (a *goError) Error() string {
	return fmt.Sprintf("%s at:\n%s", a.BaseError, string(a.Stack))
}
