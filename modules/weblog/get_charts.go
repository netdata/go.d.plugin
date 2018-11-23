package weblog

//
//import (
//	"github.com/netdata/go.d.plugin/pkg/charts"
//	"github.com/netdata/go.d.plugin/pkg/utils"
//)
//
//type (
//	Charts = charts.Charts
//	Dim    = charts.Dim
//)
//
//func (w *WebLog) GetCharts() *charts.Charts {
//	n := utils.StringSlice(w.parser.SubexpNames())
//	ch := Charts{}
//
//	ch.Add(&chartRespStatuses)
//	ch.Add(&chartRespCodes)
//
//	if w.DoCodesDetail && w.DoCodesAggregate {
//		ch.Add(&chartRespCodesDetailed)
//	}
//
//	if w.DoCodesDetail && !w.DoCodesAggregate {
//		for _, chart := range chartRespCodesDetailedPerFam() {
//			ch.Add(&chart)
//		}
//	}
//
//	if n.Include(keyBytesSent) || n.Include(keyRespLen) {
//		ch.Add(&chartBandwidth)
//	}
//
//	if n.Include(keyRequest) && w.urlCat.exist() {
//		ch.Add(chartReqPerURL.Copy())
//		for _, v := range w.urlCat.items {
//			ch.Get(chartReqPerURL.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
//			w.data[v.id] = 0
//		}
//		w.data[w.urlCat.other] = 0
//
//	}
//
//	if (n.Include(keyRequest) && w.urlCat.exist()) && w.DoChartURLCat {
//		for _, v := range w.urlCat.items {
//			for _, chart := range chartPerCategoryStats(v.id) {
//				ch.Add(&chart)
//				for _, d := range chart.Dims {
//					w.data[d.ID] = 0
//				}
//			}
//		}
//	}
//
//	if n.Include(keyUserDefined) && w.userCat.exist() {
//		ch.Add(chartReqPerUserDef.Copy())
//		for _, v := range w.userCat.items {
//			ch.Get(chartReqPerUserDef.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
//			w.data[v.id] = 0
//		}
//		w.data[w.userCat.other] = 0
//	}
//
//	if n.Include(keyRespTime) {
//		ch.Add(&chartRespTime)
//	}
//
//	if n.Include(keyRespTime) && len(w.RawHistogram) != 0 {
//		ch.Add(chartRespTimeHist.Copy())
//		for _, v := range w.histograms.get(keyRespTimeHist) {
//			ch.Get(chartRespTimeHist.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
//		}
//	}
//
//	if n.Include(keyRespTimeUpstream) {
//		ch.Add(&chartRespTimeUpstream)
//	}
//
//	if n.Include(keyRespTimeUpstream) && len(w.RawHistogram) != 0 {
//		ch.Add(chartRespTimeUpstreamHist.Copy())
//		for _, v := range w.histograms.get(keyRespTimeUpstreamHist) {
//			ch.Get(chartRespTimeUpstreamHist.ID).AddDim(&Dim{ID: v.id, Name: v.name, Algo: charts.Incremental})
//		}
//	}
//
//	if n.Include(keyRequest) {
//		ch.Add(&chartReqPerHTTPMethod)
//	}
//
//	if n.Include(keyRequest) {
//		ch.Add(&chartReqPerHTTPVer)
//	}
//
//	if n.Include(keyAddress) {
//		ch.Add(&chartReqPerIPProto)
//		ch.Add(&chartClientsCurr)
//		if w.DoClientsAll {
//			ch.Add(&chartClientsAll)
//		}
//	}
//
//	c := charts.NewCharts(ch...)
//	w.Charts = c
//	return c
//}
