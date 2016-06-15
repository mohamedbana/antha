// addresponsedata.go
package plot

import (
	"github.com/sajari/regression"
)

func Rsquared(xname string, xvalues []float64, yname string, yvalues []float64) (rsquared float64) {

	var r regression.Regression
	r.SetObservedName(yname)
	r.SetVarName(0, xname)

	for i, _ := range xvalues {
		r.AddDataPoint(regression.DataPoint{Observed: yvalues[i], Variables: []float64{xvalues[i]}})
		//r.AddDataPoint(regression.DataPoint{Observed: ControlCurvePoints + 1, Variables: ControlConcentrations})
	}
	r.RunLinearRegression()
	_ = r.GetRegCoeff(0)
	//c := r.GetRegCoeff(1)
	rsquared = r.Rsquared
	return
}
