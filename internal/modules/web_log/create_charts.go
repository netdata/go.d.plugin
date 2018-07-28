package web_log

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type (
	Charts    = charts.Charts
	Dim = charts.Dim
)

func (w *WebLog) CreateCharts() {
	n := utils.StringSlice(w.parser.SubexpNames())
	ch := Charts{}

	ch.AddChart(&chartRespStatuses)
	ch.AddChart(&chartRespCodes)

	if w.DoCodesDetail && w.DoCodesAggregate {
		ch.AddChart(&chartRespCodesDetailed)
	}

	if w.DoCodesDetail && !w.DoCodesAggregate {
		for _, chart := range chartRespCodesDetailedPerFam() {
			ch.AddChart(&chart)
		}
	}

	if n.Include(keyBytesSent) || n.Include(keyRespLen) {
		ch.AddChart(&chartBandwidth)
	}

	if (n.Include(keyRequest) || n.Include(keyURL)) && w.urlCat.exist() {
		ch.AddChart(chartReqPerURL.Copy())
		for _, v := range w.urlCat.items {
			ch.GetChart(chartReqPerURL.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.urlCat.other] = 0

	}

	if ((n.Include(keyRequest) || n.Include(keyURL)) && w.urlCat.exist()) && w.DoChartURLCat {
		for _, v := range w.urlCat.items {
			for _, chart := range chartPerCategoryStats(v.id) {
				ch.AddChart(&chart)
				for _, d := range chart.Dims {
					w.data[d.ID] = 0
				}
			}
		}
	}

	if n.Include(keyUserDefined) && w.userCat.exist() {
		ch.AddChart(chartReqPerUserDef.Copy())
		for _, v := range w.userCat.items {
			ch.GetChart(chartReqPerUserDef.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
			w.data[v.id] = 0
		}
		w.data[w.userCat.other] = 0
	}

	if n.Include(keyRespTime) {
		ch.AddChart(&chartRespTime)
	}

	if n.Include(keyRespTime) && len(w.RawHistogram) != 0 {
		ch.AddChart(chartRespTimeHist.Copy())
		for _, v := range w.histograms.get(keyRespTimeHist) {
			ch.GetChart(chartRespTimeHist.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
		}
	}

	if n.Include(keyRespTimeUpstream) {
		ch.AddChart(&chartRespTimeUpstream)
	}

	if n.Include(keyRespTimeUpstream) && len(w.RawHistogram) != 0 {
		ch.AddChart(chartRespTimeUpstreamHist.Copy())
		for _, v := range w.histograms.get(keyRespTimeUpstreamHist) {
			ch.GetChart(chartRespTimeUpstreamHist.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
		}
	}

	if n.Include(keyRequest) || n.Include(keyHTTPMethod) {
		ch.AddChart(&chartReqPerHTTPMethod)
	}

	if n.Include(keyRequest) || n.Include(keyHTTPVer) {
		ch.AddChart(&chartReqPerHTTPVer)
	}

	if n.Include(keyAddress) {
		ch.AddChart(&chartReqPerIPProto)
		ch.AddChart(&chartClientsCurr)
		if w.DoClientsAll {
			ch.AddChart(&chartClientsAll)
		}
	}

	w.AddChart(ch...)
}
