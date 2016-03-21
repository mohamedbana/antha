package pretty

import (
	"fmt"
	"io"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/target"
)

func Run(out io.Writer, result *execute.Result, t *target.Target) error {
	if len(t.Runners()) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(out, "== Running Workflow:\n"); err != nil {
		return err
	}
	for _, inst := range result.Insts {
		rinst, ok := inst.(target.RunInst)
		if !ok {
			continue
		}
		files := rinst.Data()
		rs := t.CanRun(files.Type)
		if len(rs) == 0 {
			return fmt.Errorf("no device to run type %s", files.Type)
		}
		s, err := summarize(rinst)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "    * %s", s); err != nil {
			return err
		}
		if err := rs[0].Run(files); err != nil {
			fmt.Fprintf(out, " [FAIL]\n")
			return err
		}
		if _, err := fmt.Fprintf(out, " [OK]\n"); err != nil {
			return err
		}
	}
	return nil
}
