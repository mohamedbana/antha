package inject

import (
	"errors"
	"golang.org/x/net/context"
	"sync"
)

var alreadyAdded = errors.New("already added")

type Registry struct {
	lock     sync.Mutex
	parent   context.Context
	registry map[Name]Runner
	// XXX: add remote registries...
}

// Unique identifier for Runner
type Name struct {
	Host string
	Repo string
	Tag  string
}

// Query for a Runner
type NameQuery struct {
	Repo string
	Tag  string
}

func (a *Registry) Add(name Name, runner Runner) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.registry == nil {
		a.registry = make(map[Name]Runner)
	}
	if r := a.registry[name]; r != nil {
		return alreadyAdded
	}
	a.registry[name] = runner
	return nil
}

// XXX: lift more complicated search/typing logic out
func (a *Registry) Find(query NameQuery) ([]Runner, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	name := Name{Repo: query.Repo, Tag: query.Tag}
	if r := a.registry[name]; r == nil {
		return nil, nil
	} else {
		return []Runner{r}, nil
	}
}
