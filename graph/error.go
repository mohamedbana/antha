package graph

import (
	"errors"
)

var (
	TraversalDone = errors.New("traversal done") // Predefined error for expected early termination of a traversal
	NextNode      = errors.New("next node")      // Predefined error for continuing a traversal
)
