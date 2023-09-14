// SPDX-License-Identifier: GPL-3.0-or-later

package netdataapi

import (
	"fmt"
	"io"
	"strings"
)

type (
	// API implements Netdata external plugins API.
	// https://learn.netdata.cloud/docs/agent/collectors/plugins.d#the-output-of-the-plugin
	API struct {
		io.Writer
	}
)

func New(w io.Writer) *API { return &API{w} }

// CHART  creates or update a chart.
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

// DIMENSION adds or update a dimension to the chart just created.
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

// CLABEL adds or update a label to the chart.
func (a *API) CLABEL(key, value string, source int) error {
	_, err := fmt.Fprintf(a, "CLABEL '%s' '%s' '%d'\n", key, value, source)
	return err
}

// CLABELCOMMIT adds labels to the chart. Should be called after one or more CLABEL.
func (a *API) CLABELCOMMIT() error {
	_, err := fmt.Fprint(a, "CLABEL_COMMIT\n")
	return err
}

// BEGIN initializes data collection for a chart.
func (a *API) BEGIN(typeID string, ID string, msSince int) (err error) {
	if msSince > 0 {
		_, err = fmt.Fprintf(a, "BEGIN '%s.%s' %d\n", typeID, ID, msSince)
	} else {
		_, err = fmt.Fprintf(a, "BEGIN '%s.%s'\n", typeID, ID)
	}
	return err
}

// SET sets the value of a dimension for the initialized chart.
func (a *API) SET(ID string, value int64) error {
	_, err := fmt.Fprintf(a, "SET '%s' = %d\n", ID, value)
	return err
}

// SETEMPTY sets the empty value of a dimension for the initialized chart.
func (a *API) SETEMPTY(ID string) error {
	_, err := fmt.Fprintf(a, "SET '%s' = \n", ID)
	return err
}

// VARIABLE sets the value of a CHART scope variable for the initialized chart.
func (a *API) VARIABLE(ID string, value int64) error {
	_, err := fmt.Fprintf(a, "VARIABLE CHART '%s' = %d\n", ID, value)
	return err
}

// END completes data collection for the initialized chart.
func (a *API) END() error {
	_, err := fmt.Fprintf(a, "END\n\n")
	return err
}

// FLUSH ignores the last collected values.
func (a *API) FLUSH() error {
	_, err := fmt.Fprintf(a, "FLUSH\n")
	return err
}

// DISABLE disables this plugin. This will prevent Netdata from restarting the plugin.
func (a *API) DISABLE() error {
	_, err := fmt.Fprintf(a, "DISABLE\n")
	return err
}

// EMPTYLINE writes an empty line.
func (a *API) EMPTYLINE() error {
	_, err := fmt.Fprintf(a, "\n")
	return err
}

func (a *API) HOSTINFO(guid, hostname string, labels map[string]string) error {
	if err := a.HOSTDEFINE(guid, hostname); err != nil {
		return err
	}
	for k, v := range labels {
		if err := a.HOSTLABEL(k, v); err != nil {
			return err
		}
	}
	return a.HOSTDEFINEEND()
}

func (a *API) HOSTDEFINE(guid, hostname string) error {
	_, err := fmt.Fprintf(a, "HOST_DEFINE '%s' '%s'\n", guid, hostname)
	return err
}

func (a *API) HOSTLABEL(name, value string) error {
	_, err := fmt.Fprintf(a, "HOST_LABEL '%s' '%s'\n", name, value)
	return err
}

func (a *API) HOSTDEFINEEND() error {
	_, err := fmt.Fprintf(a, "HOST_DEFINE_END\n\n")
	return err
}

func (a *API) HOST(guid string) error {
	_, err := fmt.Fprintf(a, "HOST '%s'\n\n", guid)
	return err
}

func (a *API) DynCfgEnable(pluginName string) error {
	_, err := fmt.Fprintf(a, "DYNCFG_ENABLE '%s'\n\n", pluginName)
	return err
}

func (a *API) DyncCfgRegisterModule(moduleName string) error {
	_, err := fmt.Fprintf(a, "DYNCFG_REGISTER_MODULE '%s' job_array\n\n", moduleName)
	return err
}

func (a *API) DynCfgRegisterJob(moduleName, jobName, jobType string) error {
	_, err := fmt.Fprintf(a, "DYNCFG_REGISTER_JOB '%s' '%s' '%s' 0\n\n", moduleName, jobName, jobType)
	return err
}

func (a *API) DynCfgReportJobStatus(moduleName, jobName, status, reason string) error {
	_, err := fmt.Fprintf(a, "REPORT_JOB_STATUS '%s' '%s' '%s' 0 '%s'\n\n", moduleName, jobName, status, reason)
	return err
}

func (a *API) FunctionResultSuccess(uid, contentType, payload string) error {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("FUNCTION_RESULT_BEGIN %s 0 %s 0\n", uid, contentType))
	s.WriteString(payload + "\n")
	s.WriteString("FUNCTION_RESULT_END\n\n")

	_, err := fmt.Fprintf(a, s.String())
	return err
}

func (a *API) FunctionResultReject(uid, contentType, payload string) error {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("FUNCTION_RESULT_BEGIN %s 1 %s 5\n", uid, contentType))
	s.WriteString(payload + "\n")
	s.WriteString("FUNCTION_RESULT_END\n\n")

	_, err := fmt.Fprintf(a, s.String())
	return err
}
