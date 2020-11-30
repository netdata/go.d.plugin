package nginxvts

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (nv *NginxVts) collect() (map[string]int64, error) {
	collected := make(map[string]interface{})

	ms, err := nv.apiClient.getVtsStatus()
	if err != nil {
		return nil, err
	}

	nv.addMainCharts(ms, collected)
	nv.addSharedZonesCharts(ms, collected)
	nv.addServerZonesCharts(ms, collected)
	nv.addUpstreamZonesCharts(ms, collected)
	nv.addFilterZonesCharts(ms, collected)
	nv.addCacheZonesCharts(ms, collected)

	return stm.ToMap(collected), nil
}

func (nv *NginxVts) addMainCharts(stat *vtsStatus, collected map[string]interface{}) {
	collected["loadmsec"] = stat.LoadMsec
	collected["nowmsec"] = stat.NowMsec
	collected["connections"] = stat.Connections
}

func (nv *NginxVts) addSharedZonesCharts(stat *vtsStatus, collected map[string]interface{}) {
	charts := nginxVtsSharedZonesChart.Copy()
	_ = nv.charts.Add(*charts...)
	collected["sharedzones"] = stat.SharedZones
}

func (nv *NginxVts) addServerZonesCharts(stat *vtsStatus, collected map[string]interface{}) {
	if !stat.hasServerZones() {
		return
	}

	for server := range stat.ServerZones {
		charts := nginxVtsServerZonesCharts.Copy()
		for _, chart := range *charts {
			chart.ID = fmt.Sprintf(chart.ID, server)
			chart.Fam = server
			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, server)
			}
		}
		_ = nv.charts.Add(*charts...)
	}
	collected["serverzones"] = stat.ServerZones
}

func (nv *NginxVts) addUpstreamZonesCharts(stat *vtsStatus, collected map[string]interface{}) {
	if !stat.hasUpstreamZones() {
		return
	}

	upstreamMap := make(map[string]Upstream)

	for upstreamGrp, upstreamList := range stat.UpstreamZones {
		for _, upstream := range upstreamList {
			// Merge upstream group name and upstream server as new key
			mergedKey := fmt.Sprintf("%s_%s", upstreamGrp, upstream.Server)
			upstreamMap[mergedKey] = upstream

			charts := nginxVtsUpstreamZonesCharts.Copy()
			for _, chart := range *charts {
				chart.ID = fmt.Sprintf(chart.ID, mergedKey)
				chart.Fam = upstream.Server
				for _, dim := range chart.Dims {
					dim.ID = fmt.Sprintf(dim.ID, mergedKey)
				}
			}
			_ = nv.charts.Add(*charts...)
		}
	}
	collected["upstreamzones"] = upstreamMap
}

func (nv *NginxVts) addFilterZonesCharts(stat *vtsStatus, collected map[string]interface{}) {
	if !stat.hasFilterZones() {
		return
	}

	filterMap := make(map[string]Server)

	for filter, serverMap := range stat.FilterZones {
		for group, upstream := range serverMap {
			mergedKey := fmt.Sprintf("%s_%s", filter, group)
			filterMap[mergedKey] = upstream

			charts := nginxVtsFilterZonesCharts.Copy()
			for _, chart := range *charts {
				chart.ID = fmt.Sprintf(chart.ID, mergedKey)
				chart.Fam = filter
				for _, dim := range chart.Dims {
					dim.ID = fmt.Sprintf(dim.ID, mergedKey)
				}
			}
			_ = nv.charts.Add(*charts...)
		}
	}
	collected["filterzones"] = filterMap
}

func (nv *NginxVts) addCacheZonesCharts(stat *vtsStatus, collected map[string]interface{}) {
	if !stat.hasCacheZones() {
		return
	}

	for cache := range stat.CacheZones {
		charts := nginxVtsCacheZonesCharts.Copy()
		for _, chart := range *charts {
			chart.ID = fmt.Sprintf(chart.ID, cache)
			chart.Fam = cache
			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, cache)
			}
		}
		_ = nv.charts.Add(*charts...)

	}
	collected["cachezones"] = stat.CacheZones
}
