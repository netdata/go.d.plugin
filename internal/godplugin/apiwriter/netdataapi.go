package apiwriter

import (
	"fmt"
	"io"
	"os"
)

type (
	// APIWriter write native netdata plugin API
	// https://github.com/firehol/netdata/wiki/External-Plugins#native-netdata-plugin-api
	APIWriter struct {
		// Out write to
		Out io.Writer
	}
)

// Defaults is the default API writer, which write to stdout.
var Defaults = &APIWriter{os.Stdout}

// Chart defines a new chart.
func (w *APIWriter) Chart(
	typeName string,
	ID string,
	name string,
	title string,
	units string,
	family string,
	context string,
	chartType string,
	priority int,
	updateEvery int,
	options string,
	module string) {
	fmt.Fprintf(w.Out, "CHART %s.%s '%s' '%s' '%s' '%s' '%s' %d %d %s go.d '%s'\n",
		typeName, ID, title, units, family, context, chartType, priority, updateEvery, options, module)
}

//Dimension defines a new dimension for the chart
func (w *APIWriter) Dimension(
	ID string,
	name string,
	algorithm string,
	multiplier int,
	divisor int,
	hidden string) {
	fmt.Fprintf(w.Out, "DIMENSION '%s' '%s' '%s' %d %d %s\n",
		ID, name, algorithm, multiplier, divisor, hidden)
}

// Begin initialize data collection for a chart
func (w *APIWriter) Begin(typeName string, ID string, msSince int) {
	if msSince > 0 {
		fmt.Fprintf(w.Out, "NEGIN %s.%s %d\n", typeName, ID, msSince)
	} else {
		fmt.Fprintf(w.Out, "NEGIN %s.%s\n", typeName, ID)
	}
}

// Set set the value of a dimension for the initialized chart
func (w *APIWriter) Set(ID string, value int64) {
	fmt.Fprintf(w.Out, "SET %s = %d\n", ID, value)
}

// End complete data collection for the initialized chart
func (w *APIWriter) End() {
	fmt.Fprintln(w.Out, "END")
}

// Flush ignore the last collected values
func (w *APIWriter) Flush() {
	fmt.Fprintln(w.Out, "FLUSH")
}

// Disable disable this plugin.
func (w *APIWriter) Disable() {
	fmt.Fprintln(w.Out, "DISABLE")
}
