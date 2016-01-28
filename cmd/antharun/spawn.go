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
	runCmd := func(prog string, args ...string) error {
		c := exec.Command(prog, args...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

	if err := runCmd("go", "get", gopackage); err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "server")
	if err != nil {
		return nil, err
	}
	s := &spawned{}
	s.closers = append(s.closers, func() error { return os.Remove(f.Name()) })

	if err := runCmd("go", "build", "-o", f.Name(), gopackage); err != nil {
		return s, err
	}

	prefix := ansi.Color("server", "red:white") + " "
	w1 := NewWriter(colorable.NewColorableStdout(), prefix)
	w2 := NewWriter(colorable.NewColorableStderr(), prefix)

	cmd := exec.Command(f.Name(), "-port", strconv.Itoa(port))
	cmd.Stdout = w1
	cmd.Stderr = w2

	s.Command = cmd
	s.URI = fmt.Sprintf("localhost:%d", port)
	s.closers = append(s.closers, func() error { return w1.Flush() })
	s.closers = append(s.closers, func() error { return w2.Flush() })
	return s, nil
}
