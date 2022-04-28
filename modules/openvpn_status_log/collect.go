package openvpn_status

import (
	"fmt"
	"time"
)

type clientInfo struct {
	CommonName     string
	BytesReceived  int
	BytesSent      int
	ConnectedSince int64
}

func (o *OpenVPNStatusLog) collect() (map[string]int64, error) {
	var err error

	mx := make(map[string]int64)

	clients, err := parseStatusLog(o.StatusPath)
	if err != nil {
		o.Errorf("%v", err)
		return nil, err
	}
	collectTotalStats(mx, clients)
	if o.perUserMatcher != nil {
		o.collectUsers(mx, clients)
	}

	return mx, nil
}

func collectTotalStats(mx map[string]int64, clients []clientInfo) {
	bytesIn := 0
	bytesOut := 0
	for _, c := range clients {
		bytesIn += c.BytesReceived
		bytesOut += c.BytesSent
	}
	mx["clients"] = int64(len(clients))
	mx["bytes_in"] = int64(bytesIn)
	mx["bytes_out"] = int64(bytesOut)
}

func (o *OpenVPNStatusLog) collectUsers(mx map[string]int64, clients []clientInfo) {
	now := time.Now().Unix()

	for _, user := range clients {
		name := user.CommonName
		if !o.perUserMatcher.MatchString(name) {
			continue
		}
		if !o.collectedUsers[name] {
			o.collectedUsers[name] = true
			if err := o.addUserCharts(name); err != nil {
				o.Warning(err)
			}
		}
		mx[name+"_bytes_in"] = int64(user.BytesReceived)
		mx[name+"_bytes_out"] = int64(user.BytesSent)
		mx[name+"_connection_time"] = now - user.ConnectedSince
	}
}

func (o *OpenVPNStatusLog) addUserCharts(userName string) error {
	cs := userCharts.Copy()

	for _, chart := range *cs {
		chart.ID = fmt.Sprintf(chart.ID, userName)
		chart.Fam = fmt.Sprintf(chart.Fam, userName)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, userName)
		}
		chart.MarkNotCreated()
	}
	return o.charts.Add(*cs...)
}
