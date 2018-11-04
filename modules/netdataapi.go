package modules

import (
	"fmt"
	"io"
)

type (
	// apiWriter write native netdata plugin API
	// https://github.com/firehol/netdata/wiki/External-Plugins#native-netdata-plugin-api
	apiWriter struct {
		// Out write to
		io.Writer
	}
)

// chart defines a new chart.
func (w *apiWriter) chart(
	typeName string,
	ID string,
	name string,
	title string,
	units string,
	family string,
	context string,
	chartType chartType,
	priority int,
	updateEvery int,
	options Opts,
	module string) error {
	_, err := fmt.Fprintf(w, "CHART %s.%s '%s' '%s' '%s' '%s' '%s' '%s' %d %d %s go.d '%s'\n",
		typeName, ID, name, title, units, family, context, chartType, priority, updateEvery, options, module)
	return err
}

//dimension defines a new dimension for the chart
func (w *apiWriter) dimension(
	ID string,
	name string,
	algorithm string,
	multiplier int,
	divisor int,
	hidden string) error {
	_, err := fmt.Fprintf(w, "DIMENSION '%s' '%s' '%s' %d %d %s\n",
		ID, name, algorithm, multiplier, divisor, hidden)
	return err
}

// begin initialize data collection for a chart
func (w *apiWriter) begin(typeName string, ID string, msSince int) error {
	var err error
	if msSince > 0 {
		_, err = fmt.Fprintf(w, "BEGIN %s.%s %d\n", typeName, ID, msSince)
	} else {
		_, err = fmt.Fprintf(w, "BEGIN %s.%s\n", typeName, ID)
	}
	return err
}

// set set the value of a dimension for the initialized chart
func (w *apiWriter) set(ID string, value int64) error {
	_, err := fmt.Fprintf(w, "SET %s = %d\n", ID, value)
	return err
}

// end complete data collection for the initialized chart
func (w *apiWriter) end() error {
	_, err := fmt.Fprintf(w, "END\n")
	return err
}

// flush ignore the last collected values
func (w *apiWriter) flush() error {
	_, err := fmt.Fprintf(w, "FLUSH\n")
	return err
}
