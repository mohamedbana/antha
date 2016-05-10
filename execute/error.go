package execute

import (
	"fmt"

	"golang.org/x/net/context"
)

type Error struct {
	Message string
}

func (a *Error) Error() string {
	return a.Message
}

// Report an execution error. Does not return
func Errorf(ctx context.Context, format string, args ...interface{}) {
	panic(&Error{Message: fmt.Sprintf(format, args...)})
}
