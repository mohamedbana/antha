package main

import (
	"fmt"
	"github.com/antha-lang/antha/internal/github.com/mattn/go-colorable"
	"github.com/antha-lang/antha/internal/github.com/mgutz/ansi"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

type fn func() error
type spawned struct {
	Command *exec.Cmd
	URI     string
	closers []fn
}

func (a *spawned) Close() (err error) {
	for _, c := range a.closers {
		if e := c(); e != nil {
			err = e
		}
	}
	return
}

func spawn(gopackage string, port int) (*spawned, error) {
	if f, err := ioutil.TempFile("", "server"); err != nil {
		return nil, err
	} else if err := exec.Command("go", "build", "-o", f.Name(), gopackage).Run(); err != nil {
		return nil, err
	} else {
		prefix := ansi.Color("server", "red:white") + " "
		cmd := exec.Command(f.Name(), "-port", strconv.Itoa(port))
		w1 := NewWriter(colorable.NewColorableStdout(), prefix)
		w2 := NewWriter(colorable.NewColorableStderr(), prefix)
		cmd.Stdout = w1
		cmd.Stderr = w2
		s := &spawned{
			Command: cmd,
			URI:     fmt.Sprintf("localhost:%d", port),
		}
		s.closers = append(s.closers, func() error { return os.Remove(f.Name()) })
		s.closers = append(s.closers, func() error { return w1.Flush() })
		s.closers = append(s.closers, func() error { return w2.Flush() })
		return s, nil
	}
}
