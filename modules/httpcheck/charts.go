package httpcheck

//
//import (
//	"github.com/netdata/go.d.plugin/pkg/charts"
//)
//
//type (
//	Charts = charts.Charts
//	Opts   = charts.Opts
//	Dims   = charts.Dims
//)
//
//var uCharts = Charts{
//	{
//		ID:   "response_time",
//		Opts: Opts{Title: "HTTP Response Time", Units: "ms", Fam: "response"},
//		Dims: Dims{
//			{ID: "response_time", Name: "time", Div: 1000000},
//		},
//	},
//	{
//		ID:   "response_length",
//		Opts: Opts{Title: "HTTP Response Body Length", Units: "characters", Fam: "response"},
//		Dims: Dims{
//			{ID: "response_length", Name: "length", Div: 1000000},
//		},
//	},
//	{
//		ID:   "response_status",
//		Opts: Opts{Title: "HTTP Response Status", Units: "boolean", Fam: "status"},
//		Dims: Dims{
//			{ID: "success"},
//			{ID: "failed"},
//			{ID: "timeout"},
//		},
//	},
//	{
//		ID:   "response_check_status",
//		Opts: Opts{Title: "HTTP Response Check Status", Units: "boolean", Fam: "status"},
//		Dims: Dims{
//			{ID: "bad_status", Name: "bad status"},
//		},
//	},
//	{
//		ID:   "response_check_content",
//		Opts: Opts{Title: "HTTP Response Check Content", Units: "boolean", Fam: "status"},
//		Dims: Dims{
//			{ID: "bad_content", Name: "bad content"},
//		},
//	},
//}
