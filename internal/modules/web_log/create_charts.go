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

func (w *WebLog) CreateCharts() {
	n := utils.StringSlice(w.parser.SubexpNames())
	ch := new(Charts)

	ch.AddChart(charts.RespStatuses)
	ch.AddChart(charts.RespCodes)

	if w.DoCodesDetail && w.DoCodesAggregate {
		ch.AddChart(charts.RespCodesDetailed)
	}

	if w.DoCodesDetail && !w.DoCodesAggregate {
		for _, chart := range charts.RespCodesDetailedPerFam() {
			ch.AddChart(chart)
		}
	}

	if n.Include(keyBytesSent) || n.Include(keyRespLen) {
		ch.AddChart(charts.Bandwidth)
	}

	if (n.Include(keyRequest) || n.Include(keyURL)) && w.urlCat.exist() {
		ch.AddChart(charts.ReqPerURL)
		for _, v := range w.urlCat.items {
			ch.GetChartByID(charts.ReqPerURL.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.urlCat.other] = 0

	}

	if ((n.Include(keyRequest) || n.Include(keyURL)) && w.urlCat.exist()) && w.DoChartURLCat {
		for _, v := range w.urlCat.items {
			for _, chart := range charts.PerCategoryStats(v.id) {
				ch.AddChart(chart)
				for _, d := range chart.Dimensions {
					w.data[d.ID()] = 0
				}
			}
		}
	}

	if n.Include(keyUserDefined) && w.userCat.exist() {
		ch.AddChart(charts.ReqPerUserDef)
		for _, v := range w.userCat.items {
			ch.GetChartByID(charts.ReqPerUserDef.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.userCat.other] = 0
	}

	if n.Include(keyRespTime) {
		ch.AddChart(charts.RespTime)
	}

	if n.Include(keyRespTime) && len(w.RawHistogram) != 0 {
		ch.AddChart(charts.RespTimeHist)
		for _, v := range w.histograms.get(keyRespTimeHist) {
			ch.GetChartByID(charts.RespTimeHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
		}
	}

	if n.Include(keyRespTimeUpstream) {
		ch.AddChart(charts.RespTimeUpstream)
	}

	if n.Include(keyRespTimeUpstream) && len(w.RawHistogram) != 0 {
		ch.AddChart(charts.RespTimeUpstreamHist)
		for _, v := range w.histograms.get(keyRespTimeUpstreamHist) {
			ch.GetChartByID(charts.RespTimeUpstreamHist.ID).AddDim(Dimension{v.id, v.name, raw.Incremental})
		}
	}

	if n.Include(keyRequest) || n.Include(keyHTTPMethod) {
		ch.AddChart(charts.ReqPerHTTPMethod)
	}

	if n.Include(keyRequest) || n.Include(keyHTTPVer) {
		ch.AddChart(charts.ReqPerHTTPVer)
	}

	if n.Include(keyAddress) {
		ch.AddChart(charts.ReqPerIPProto)
		ch.AddChart(charts.ClientsCurr)
		if w.DoClientsAll {
			ch.AddChart(charts.ClientsAll)
		}
	}

	w.AddMany(ch)
}
