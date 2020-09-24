package netdataapi

import (
	"fmt"
	"io"
)

type (
	// API implements Netdata external plugins API.
	// https://learn.netdata.cloud/docs/agent/collectors/plugins.d#the-output-of-the-plugin
	API struct {
		io.Writer
	}
)

func New(w io.Writer) *API { return &API{w} }

// CHART  create or update a chart.
func (a *API) CHART(
	typeID string,
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
	plugin string,
	module string) error {
	_, err := fmt.Fprintf(a, "CHART '%s.%s' '%s' '%s' '%s' '%s' '%s' '%s' '%d' '%d' '%s' '%s' '%s'\n",
		typeID, ID, name, title, units, family, context, chartType, priority, updateEvery, options, plugin, module)
	return err
}

// DIMENSION add or update a dimension to the chart just created.
func (a *API) DIMENSION(
	ID string,
	name string,
	algorithm string,
	multiplier int,
	divisor int,
	options string) error {
	_, err := fmt.Fprintf(a, "DIMENSION '%s' '%s' '%s' '%d' '%d' '%s'\n",
		ID, name, algorithm, multiplier, divisor, options)
	return err
}

// BEGIN initialize data collection for a chart.
func (a *API) BEGIN(typeID string, ID string, msSince int) (err error) {
	if msSince > 0 {
		_, err = fmt.Fprintf(a, "BEGIN '%s.%s' %d\n", typeID, ID, msSince)
	} else {
		_, err = fmt.Fprintf(a, "BEGIN '%s.%s'\n", typeID, ID)
	}
	return err
}

// SET set the value of a dimension for the initialized chart.
func (a *API) SET(ID string, value int64) error {
	_, err := fmt.Fprintf(a, "SET '%s' = %d\n", ID, value)
	return err
}

// SETEMPTY set the empty value of a dimension for the initialized chart.
func (a *API) SETEMPTY(ID string) error {
	_, err := fmt.Fprintf(a, "SET '%s' = \n", ID)
	return err
}

// VARIABLE set the value of a CHART scope variable for the initialized chart.
func (a *API) VARIABLE(ID string, value int64) error {
	_, err := fmt.Fprintf(a, "VARIABLE CHART '%s' = %d\n", ID, value)
	return err
}

// END complete data collection for the initialized chart.
func (a *API) END() error {
	_, err := fmt.Fprintf(a, "END\n\n")
	return err
}

// FLUSH ignore the last collected values.
func (a *API) FLUSH() error {
	_, err := fmt.Fprintf(a, "FLUSH\n")
	return err
}

// DISABLE disable this plugin. This will prevent Netdata from restarting the plugin.
func (a *API) DISABLE() error {
	_, err := fmt.Fprintf(a, "DISABLE\n")
	return err
}

// EMPTYLINE write an empty line.
func (a *API) EMPTYLINE() error {
	_, err := fmt.Fprintf(a, "\n")
	return err
}
