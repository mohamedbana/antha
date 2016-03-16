package target

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	driver "github.com/antha-lang/antha/driver/pb"
)

type Runner struct {
	client driver.RunnerClient
}

func NewRunner(c driver.RunnerClient) *Runner {
	return &Runner{c}
}

func loadFile(name string, tarball []byte) ([]byte, error) {
	buf := bytes.NewBuffer(tarball)
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	r := tar.NewReader(gr)
	var ret bytes.Buffer
	for {
		h, err := r.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("%s not found", name)
		} else if err != nil {
			return nil, err
		}
		out := ioutil.Discard

		if h.Name == name {
			out = &ret
		}
		if _, err := io.Copy(out, r); err != nil {
			return nil, err
		}

		if h.Name == name {
			return ret.Bytes(), nil
		}
	}
}

func (a *Runner) Run(files Files) error {
	bs, err := loadFile("input", files.Tarball)
	if err != nil {
		return err
	}
	ctx := context.Background()
	reply, err := a.client.Run(ctx, &driver.RunRequest{
		Type: files.Type,
		Data: bs,
	})
	if err != nil {
		return err
	}

	// Proof of concept
	var errors []string
	for range time.Tick(5 * time.Second) {
		msgs, err := a.client.Messages(ctx, &driver.MessagesRequest{
			Id: reply.Id,
		})
		if err != nil {
			return err
		}
		for _, m := range msgs.Values {
			if m.Code == "error" {
				errors = append(errors, string(m.Data))
			}
			if m.Code == "fatal" {
				return fmt.Errorf("error running protocol: %s", strings.Join(errors, " "))
			}
			if m.Code == "stop" {
				return nil
			}
		}
	}

	return nil
}

func (a *Runner) types() ([]string, error) {
	ctx := context.Background()
	reply, err := a.client.SupportedRunTypes(ctx, &driver.SupportedRunTypesRequest{})
	if err != nil {
		return nil, err
	}
	return reply.Types, nil
}

func (a *Runner) String() string {
	return "Runner"
}

func (a *Runner) CanCompile(ast.Request) bool {
	return false
}

func (a *Runner) MoveCost(Device) int {
	return 0
}

func (a *Runner) Compile([]ast.Command) ([]Inst, error) {
	return nil, nil
}
