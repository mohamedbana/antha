package flow

import (
	"github.com/antha-lang/antha/internal/code.google.com/p/go.net/websocket"
)

type nodeHandler func(*websocket.Conn, interface{})

// Runtime is a NoFlo-compatible runtime implementing the FBP protocol
type Node struct {
	id       string
	handlers map[string]protocolHandler
	ready    chan struct{}
	done     chan struct{}
}
