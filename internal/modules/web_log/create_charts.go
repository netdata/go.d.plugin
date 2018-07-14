package web_log

import (
	"github.com/l2isbad/go.d.plugin/internal/modules/web_log/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type (
	Charts    = raw.Charts
	Dimension = raw.Dimension
)

func (w *WebLog) createCharts() {
	n := utils.StringSlice(w.regex.parser.SubexpNames())
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

	if n.Include(keyBytesSent) || n.Include(keyRespLen) {
		c.AddChart(charts.Bandwidth, true)
	}

	if n.Include(keyRespTime) {
		c.AddChart(charts.RespTime, true)

		if h := w.histograms[keyRespTimeHist]; h != nil {
			c.AddChart(charts.RespTimeHist, true)
			for _, v := range *h {
				c.GetChartByID(charts.RespTimeHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
				w.data[v.id] = 0
			}
		}
	}

	if n.Include(keyRespTimeUp) {
		c.AddChart(charts.RespTimeUpstream, true)

		if h := w.histograms[keyRespTimeUpHist]; h != nil {
			c.AddChart(charts.RespTimeUpstreamHist, true)
			for _, v := range *h {
				c.GetChartByID(charts.RespTimeUpstreamHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
				w.data[v.id] = 0
			}
		}
	}

	if n.Include(keyRequest) && w.regex.URLCat.active() {
		c.AddChart(charts.ReqPerURL, true)
		for _, v := range w.regex.URLCat.list {
			c.GetChartByID(charts.ReqPerURL.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.regex.URLCat.other()] = 0

		if w.DoChartURLCat {
			for _, v := range w.regex.URLCat.list {
				for _, chart := range charts.PerCategory(v.id) {
					c.AddChart(chart, true)
					for _, d := range chart.Dimensions {
						w.data[d.ID()] = 0
					}
				}
			}
		}
	}

	if n.Include(keyUserDefined) && w.regex.UserCat.active() {
		c.AddChart(charts.ReqPerUserDef, true)
		for _, v := range w.regex.UserCat.list {
			c.GetChartByID(charts.ReqPerUserDef.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.regex.UserCat.other()] = 0
	}

	if n.Include(keyRequest) {
		c.AddChart(charts.ReqPerHTTPMethod, true)
		c.AddChart(charts.ReqPerHTTPVer, true)
	}

	if n.Include(keyAddress) {
		c.AddChart(charts.ReqPerIPProto, true)
		c.AddChart(charts.ClientsCurr, true)
		if w.DoClientsAll {
			c.AddChart(charts.ClientsAll, true)
		}
	}

	w.AddMany(c)
}
