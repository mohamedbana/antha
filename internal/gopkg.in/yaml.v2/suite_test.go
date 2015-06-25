package yaml_test

import (
	. "github.com/antha-lang/antha/internal/gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})
