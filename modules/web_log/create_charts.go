package web_log

import (
	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules/web_log/charts"
	"github.com/l2isbad/go.d.plugin/shared"
)

type (
	Charts    = raw.Charts
	Dimension = raw.Dimension
)

func (w *WebLog) createCharts() {
	names := shared.StringSlice(w.regex.parser.SubexpNames())
	c := &Charts{}

	c.AddChart(charts.RespStatuses, true)
	c.AddChart(charts.RespCodes, true)

	if w.DoDetailCodes {
		if w.DoDetailCodesA {
			c.AddChart(charts.RespCodesDetailed, true)
		} else {
			for _, chart := range charts.RespCodesDetailedPerFam() {
				c.AddChart(chart, true)
			}
		}
	}

	if names.Include(keyBytesSent) || names.Include(keyRespLen) {
		c.AddChart(charts.Bandwidth, true)
	}

	if names.Include(keyRespTime) {
		c.AddChart(charts.RespTime, true)

		if h := w.histograms[keyRespTimeHist]; h != nil {
			c.AddChart(charts.RespTimeHist, true)
			for _, v := range *h {
				c.GetChartByID(charts.RespTimeHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
				w.data[v.id] = 0
			}
		}
	}

	if names.Include(keyRespTimeUp) {
		c.AddChart(charts.RespTimeUpstream, true)

		if h := w.histograms[keyRespTimeUpHist]; h != nil {
			c.AddChart(charts.RespTimeUpstreamHist, true)
			for _, v := range *h {
				c.GetChartByID(charts.RespTimeUpstreamHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
				w.data[v.id] = 0
			}
		}
	}

	if names.Include(keyRequest) && w.regex.URLCat.active() {
		c.AddChart(charts.ReqPerURL, true)
		for _, v := range w.regex.URLCat.list {
			c.GetChartByID(charts.ReqPerURL.ID).AddDim(Dimension{v.fullname, v.name, raw.Incremental})
			w.data[v.fullname] = 0
		}

		if w.DoChartURLCat {
			for _, v := range w.regex.URLCat.list {
				for _, chart := range charts.PerCategory(v.fullname) {
					c.AddChart(chart, true)
					for _, d := range chart.Dimensions {
						w.data[d.ID()] = 0
					}
				}
			}
		}
	}

	if names.Include(keyUserDefined) && w.regex.UserCat.active() {
		c.AddChart(charts.ReqPerUserDef, true)
		for _, v := range w.regex.UserCat.list {
			c.GetChartByID(charts.ReqPerUserDef.ID).AddDim(Dimension{v.fullname, v.name, raw.Incremental})
			w.data[v.fullname] = 0
		}
	}

	if names.Include(keyRequest) {
		c.AddChart(charts.ReqPerHTTPMethod, true)
		c.AddChart(charts.ReqPerHTTPVer, true)
	}

	if names.Include(keyAddress) {
		c.AddChart(charts.ReqPerIPProto, true)
		c.AddChart(charts.ClientsCurr, true)
		if w.DoClientsAll {
			c.AddChart(charts.ClientsAll, true)
		}
	}

	w.AddMany(c)
}
