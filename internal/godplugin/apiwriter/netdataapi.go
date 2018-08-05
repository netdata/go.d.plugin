package apiwriter

import (
	"fmt"
	"io"
)

type (
	// APIWriter write native netdata plugin API
	// https://github.com/firehol/netdata/wiki/External-Plugins#native-netdata-plugin-api
	APIWriter struct {
		// Out write to
		io.Writer
	}
)

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
	module string) error {
	_, err := fmt.Fprintf(w, "CHART %s.%s '%s' '%s' '%s' '%s' '%s' %d %d %s go.d '%s'\n",
		typeName, ID, title, units, family, context, chartType, priority, updateEvery, options, module)
	return err
}

//Dimension defines a new dimension for the chart
func (w *APIWriter) Dimension(
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

// Begin initialize data collection for a chart
func (w *APIWriter) Begin(typeName string, ID string, msSince int) error {
	var err error
	if msSince > 0 {
		_, err = fmt.Fprintf(w, "NEGIN %s.%s %d\n", typeName, ID, msSince)
	} else {
		_, err = fmt.Fprintf(w, "NEGIN %s.%s\n", typeName, ID)
	}
	return err
}

// Set set the value of a dimension for the initialized chart
func (w *APIWriter) Set(ID string, value int64) error {
	_, err := fmt.Fprintf(w, "SET %s = %d\n", ID, value)
	return err
}

// End complete data collection for the initialized chart
func (w *APIWriter) End() error {
	_, err := fmt.Fprintln(w, "END")
	return err
}

// Flush ignore the last collected values
func (w *APIWriter) Flush() error {
	_, err := fmt.Fprintln(w, "FLUSH")
	return err
}

// Disable disable this plugin.
func (w *APIWriter) Disable() error {
	_, err := fmt.Fprintln(w, "DISABLE")
	return err
}
