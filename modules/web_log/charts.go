package web_log

import (
	"github.com/l2isbad/go.d.plugin/charts/raw"
)

const (
	chartDetRespCodes  = "detailed_response_codes"
	chartHttpMethod    = "http_method"
	chartHttpVersion   = "http_version"
	chartReqPerURL     = "requests_per_url"
	chartReqPerUserDef = "requests_per_user_defined"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

var uCharts = Charts{
	Order: Order{
		"response_statuses", // fam: responses
		"response_codes",    // fam: responses
		// detailed_response_codes               // fam: responses
		// detailed_response_codes_(1xx|2xx|...) // fam: responses
		"bandwidth",              // fam: bandwidth
		"response_time",          // fam: timings
		"response_time_upstream", // fam: timings
		"requests_per_url",       // fam: urls
		// url_XXX_detailed_response_codes  //fam: url XXX
		// url_XXX_bandwidth                //fam: url XXX
		// url_XXX_response_time            //fam: url XXX
		"requests_per_user_defined", // fam: user defined
		"http_method",               // fam: http methods
		"http_version",              // fam: http versions
		"requests_per_ipproto",      // fam: ip protocols
		"clients",                   // fam: clients
		"clients_all",               // fam: clients
	},
	Definitions: Definitions{
		Chart{
			ID:      "response_statuses",
			Options: Options{"Response Statuses", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"successful_requests", "success", raw.Incremental},
				Dimension{"server_errors", "error", raw.Incremental},
				Dimension{"redirects", "redirect", raw.Incremental},
				Dimension{"bad_requests", "bad", raw.Incremental},
				Dimension{"other_requests", "other", raw.Incremental},
			},
		},
		Chart{
			ID:      "response_codes",
			Options: Options{"Response Codes", "requests/s", "responses", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"2xx", "", raw.Incremental},
				Dimension{"5xx", "", raw.Incremental},
				Dimension{"3xx", "", raw.Incremental},
				Dimension{"4xx", "", raw.Incremental},
				Dimension{"1xx", "", raw.Incremental},
				Dimension{"0xx", "", raw.Incremental},
				Dimension{"unmatched", "", raw.Incremental},
			},
		},
		Chart{
			ID:      "bandwidth",
			Options: Options{"Bandwidth", "kilobits/s", "bandwidth", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_length", "received", 8, 1000, raw.Incremental},
				Dimension{"bytes_sent", "sent", -8, 1000, raw.Incremental},
			},
		},
		Chart{
			ID:      "response_time",
			Options: Options{"Processing Time", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:      "response_time_upstream",
			Options: Options{"Processing Time Upstream", "milliseconds", "timings", "", raw.Area},
			Dimensions: Dimensions{
				Dimension{"resp_time_upstream_min", "min", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_max", "max", raw.Incremental, 1, 1000},
				Dimension{"resp_time_upstream_avg", "avg", raw.Incremental, 1, 1000},
			},
		},
		Chart{
			ID:         "requests_per_url",
			Options:    Options{"Requests Per Url", "requests/s", "urls", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:         "requests_per_user_defined",
			Options:    Options{"Requests Per User Defined Pattern", "requests/s", "user defined", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:      "http_method",
			Options: Options{"Requests Per HTTP Method", "requests/s", "http methods", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"GET", "", raw.Incremental},
			},
		},
		Chart{
			ID:         "http_version",
			Options:    Options{"Requests Per HTTP Version", "requests/s", "http versions", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:         "requests_per_ipproto",
			Options:    Options{"Requests Per IP Protocol", "requests/s", "ip protocols", "", raw.Stacked},
			Dimensions: Dimensions{},
		},
		Chart{
			ID:      "clients",
			Options: Options{"Current Poll Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_cur_ipv4", "ipv4", raw.Incremental},
				Dimension{"unique_cur_ipv6", "ipv6", raw.Incremental},
			},
		},
		Chart{
			ID:      "clients_all",
			Options: Options{"All Time Unique Client IPs", "unique ips", "clients", "", raw.Stacked},
			Dimensions: Dimensions{
				Dimension{"unique_tot_ipv4", "ipv4"},
				Dimension{"unique_tot_ipv6", "ipv6"},
			},
		},
	},
}

func (w *WebLog) addCharts() {
	c := uCharts.Copy()
	if w.DetRespCodes {
		for _, chart := range detRespCodesCharts(w.DetRespCodesA) {
			c.AddChart(chart, false)
			c.Order.InsertBefore("bandwidth", chart.ID)
		}
	}

	for _, v := range w.regex.URLCat.list {
		c.GetChartByID(chartReqPerURL).AddDim(Dimension{v.fullname, v.name, raw.Incremental})
		w.data[v.fullname] = 0
	}

	for _, v := range w.regex.UserCat.list {
		c.GetChartByID(chartReqPerUserDef).AddDim(Dimension{v.fullname, v.name, raw.Incremental})
		w.data[v.fullname] = 0
	}

	if w.ChartURLCat {
		for _, v := range w.regex.URLCat.list {
			for _, chart := range perCategoryCharts(v) {
				c.AddChart(chart, false)
				c.Order.InsertBefore(chartReqPerUserDef, chart.ID)
			}
		}
	}
	w.AddMany(c)
}

func perCategoryCharts(c *category) []raw.Chart {
	return []raw.Chart{
		raw.NewChart(
			c.fullname+"_"+chartDetRespCodes,
			Options{"Detailed Response Codes", "requests/s", c.fullname, "web_log.url_detailed_response_codes", raw.Stacked},
		),
		raw.NewChart(
			c.fullname+"_bandwidth",
			Options{"Bandwidth", "kilobits/s", c.fullname, "web_log.url_bandwidth", raw.Area},
			Dimension{c.fullname + "_resp_length", "received", raw.Incremental, 8, 1000},
			Dimension{c.fullname + "_bytes_sent", "sent", raw.Incremental, -8, 1000},
		),
		raw.NewChart(
			c.fullname+"_response_time",
			Options{"Processing Time", "milliseconds", c.fullname, "web_log.url_response_time", raw.Area},
			Dimension{c.fullname + "_resp_time_min", "min", raw.Incremental, 1, 1000},
			Dimension{c.fullname + "_resp_time_max", "max", raw.Incremental, 1, 1000},
			Dimension{c.fullname + "_resp_time_avg", "avg", raw.Incremental, 1, 1000},
		),
	}
}

func detRespCodesCharts(aggregate bool) []raw.Chart {
	if aggregate {
		return []raw.Chart{
			raw.NewChart(
				chartDetRespCodes,
				Options{"Detailed Response Codes", "requests/s", "responses", "", raw.Stacked}),
		}
	}
	return []raw.Chart{
		raw.NewChart(
			chartDetRespCodes+"_1xx",
			Options{"Detailed Response Codes 1xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			chartDetRespCodes+"_2xx",
			Options{"Detailed Response Codes 2xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			chartDetRespCodes+"_3xx",
			Options{"Detailed Response Codes 3xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			chartDetRespCodes+"_4xx",
			Options{"Detailed Response Codes 4xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			chartDetRespCodes+"_5xx",
			Options{"Detailed Response Codes 5xx", "requests/s", "responses", "", raw.Stacked},
		),
		raw.NewChart(
			chartDetRespCodes+"_other",
			Options{"Detailed Response Codes Other", "requests/s", "responses", "", raw.Stacked},
		),
	}
}
