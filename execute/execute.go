package execute

import (
	"fmt"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/workflow"
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

type Options struct {
	Id           string
	WorkflowData []byte
	ParamData    []byte
}

func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData})
	if err != nil {
		return nil, err
	}

	cd, err := setParams(parent, opt.ParamData, w)
	if err != nil {
		return nil, fmt.Errorf("cannot set initial parameters: %s", err)
	}

	ctx, done, err := newLHContext(context.WithValue(parent, theIdKey, opt.Id), cd)
	if done != nil {
		defer done()
	}
	if err != nil {
		return nil, fmt.Errorf("cannot initialize liquid handler: %s", err)
	}

	// XXX: Setup trace

	if err := w.Run(ctx); err != nil {
		return nil, err
	}

	return w, nil

	//cfData, err := ioutil.ReadFile(parametersFile)
	//if err != nil {
	//	return err
	//}

	//cf, err := util.NewConfig(cfData, wf)
	//if err != nil {
	//	return err
	//}
	//if _, ok := cf.Config["JOBID"]; !ok {
	//	cf.Config["JOBID"] = "default"
	//}
}
