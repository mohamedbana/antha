package execute

import (
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

type idKey int

const theIdKey idKey = 0

func getId(ctx context.Context) string {
	v, ok := ctx.Value(theIdKey).(string)
	if !ok {
		return ""
	}
	return v
}

func WithId(parent context.Context, id string) context.Context {
	return context.WithValue(parent, theIdKey, id)
}
