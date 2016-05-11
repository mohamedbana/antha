package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

const (
	NumConcurrent = 8
)

type TInput struct {
	WorkflowPath string
	WorkflowData []byte
	ParamPath    string
	ParamData    []byte
	Dir          string
}

type TInputs []*TInput

func (a TInputs) Len() int {
	return len(a)
}

func (a TInputs) Less(i, j int) bool {
	if a[i].WorkflowPath == a[j].WorkflowPath {
		return a[i].ParamPath < a[j].ParamPath
	}
	return a[i].WorkflowPath < a[j].WorkflowPath
}

func (a TInputs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func findInputs(basePath string) ([]*TInput, error) {
	wfiles := make(map[string][]string)
	pfiles := make(map[string][]string)
	walk := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		pabs, err := filepath.Abs(p)
		if err != nil {
			return err
		}

		dir := filepath.Dir(pabs)
		b := filepath.Base(pabs)
		if ridx := strings.LastIndex(b, "."); ridx >= 0 && strings.HasSuffix(b[:ridx], "workflow") {
			wfiles[dir] = append(wfiles[dir], pabs)
		}

		if ridx := strings.LastIndex(b, "."); ridx >= 0 && strings.HasSuffix(b[:ridx], "parameters") {
			pfiles[dir] = append(pfiles[dir], pabs)
		}
		return nil
	}

	if len(basePath) == 0 {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	if err := filepath.Walk(basePath, walk); err != nil {
		return nil, err
	}

	var inputs []*TInput
	for dir, wfs := range wfiles {
		pfs := pfiles[dir]
		switch nwfs, npfs := len(wfs), len(pfs); {
		case nwfs == 0 || npfs == 0:
			continue
		case nwfs == npfs:
			sort.Strings(wfs)
			sort.Strings(pfs)
			for idx := range wfs {
				inputs = append(inputs, &TInput{
					WorkflowPath: wfs[idx],
					ParamPath:    pfs[idx],
					Dir:          dir,
				})
			}
		case nwfs == 1:
			for idx := range pfs {
				inputs = append(inputs, &TInput{
					WorkflowPath: wfs[0],
					ParamPath:    pfs[idx],
					Dir:          dir,
				})
			}
		default:
			continue
		}
	}

	for _, input := range inputs {
		wfdata, err := ioutil.ReadFile(input.WorkflowPath)
		if err != nil {
			return nil, err
		}
		pfdata, err := ioutil.ReadFile(input.ParamPath)
		if err != nil {
			return nil, err
		}
		input.ParamData = pfdata
		input.WorkflowData = wfdata
	}

	return inputs, nil
}

func runElements(t *testing.T, ctx context.Context, inputs []*TInput) {
	tgt := target.New()
	tgt.AddDevice(human.New(human.Opt{CanMix: true, CanIncubate: true}))

	odir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range inputs {
		if len(input.Dir) != 0 {
			if err := os.Chdir(input.Dir); err != nil {
				t.Fatal(err)
			}
		}
		t.Logf("Running %q %q\n", input.WorkflowPath, input.ParamPath)
		_, err := execute.Run(ctx, execute.Opt{
			WorkflowData: input.WorkflowData,
			ParamData:    input.ParamData,
			Target:       tgt,
		})

		if err == nil {
			return
		} else if _, ok := err.(*execute.Error); ok {
			return
		} else {
			t.Errorf("error running with workflow %q with parameters %q: %s", input.WorkflowPath, input.ParamPath, err)
		}
	}

	if err := os.Chdir(odir); err != nil {
		t.Fatal(err)
	}
}

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range GetComponents() {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func getExampleInputs(t *testing.T) []*TInput {
	flag.Parse()
	args := flag.Args()
	input := "../../examples"
	if len(args) != 0 {
		input = args[0]
	}

	inputs, err := findInputs(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(inputs) == 0 {
		t.Fatalf("no tests found under path %q", input)
	}

	sort.Sort(TInputs(inputs))

	return inputs
}

// Divide l into n pieces, return indices for ith piece
func divide(i, n, l int) (int, int) {
	each := (l + n - 1) / n
	if i == n-1 {
		return i * each, l
	} else {
		return i * each, (i + 1) * each
	}
}

func TestElementsWithExampleInputs0(t *testing.T) {
	t.Parallel()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs := getExampleInputs(t)
	first, last := divide(0, 5, len(inputs))

	runElements(t, ctx, inputs[first:last])
}

func TestElementsWithExampleInputs1(t *testing.T) {
	t.Parallel()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs := getExampleInputs(t)
	first, last := divide(1, 5, len(inputs))

	runElements(t, ctx, inputs[first:last])
}

func TestElementsWithExampleInputs2(t *testing.T) {
	t.Parallel()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs := getExampleInputs(t)
	first, last := divide(2, 5, len(inputs))

	runElements(t, ctx, inputs[first:last])
}

func TestElementsWithExampleInputs3(t *testing.T) {
	t.Parallel()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs := getExampleInputs(t)
	first, last := divide(3, 5, len(inputs))

	runElements(t, ctx, inputs[first:last])
}

func TestElementsWithExampleInputs4(t *testing.T) {
	t.Parallel()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs := getExampleInputs(t)
	first, last := divide(4, 5, len(inputs))

	runElements(t, ctx, inputs[first:last])
}

func TestElementsWithDefaultInputs(t *testing.T) {
	t.Parallel()
	type Process struct {
		Component string `json:"component"`
	}
	type Workflow struct {
		Processes map[string]Process `json:"processes"`
	}
	var inputs []*TInput
	for _, c := range GetComponents() {
		wf := &Workflow{
			Processes: map[string]Process{
				"Process": {
					Component: c.Name,
				},
			},
		}
		bs, err := json.Marshal(wf)
		if err != nil {
			t.Fatal(err)
		}

		inputs = append(inputs, &TInput{
			WorkflowPath: c.Name,
			WorkflowData: bs,
		})
	}

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}
	runElements(t, ctx, inputs)
}
