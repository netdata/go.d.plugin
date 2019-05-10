package openvpn

import (
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/modules/openvpn/client"
)

func (o *OpenVPN) collect() (map[string]int64, error) {
	var err error
	if !o.apiClient.IsConnected() {
		if err = o.apiClient.Connect(); err != nil {
			return nil, err
		}
	}

	defer func() {
		// TODO: disconnect not on every error?
		if err != nil {
			_ = o.apiClient.Disconnect()
		}
	}()

	mx := make(map[string]int64)

	if err = o.collectLoadStats(mx); err != nil {
		return nil, err
	}

	if o.perUserMatcher != nil {
		if err = o.collectUsers(mx); err != nil {
			return nil, err
		}
	}

	return mx, nil
}

func (o *OpenVPN) collectLoadStats(mx map[string]int64) error {
	stats, err := o.apiClient.GetLoadStats()
	if err != nil {
		return err
	}

	mx["clients"] = stats.NumOfClients
	mx["bytes_in"] = stats.BytesIn
	mx["bytes_out"] = stats.BytesOut

	return nil
}

func (o *OpenVPN) collectUsers(mx map[string]int64) error {
	users, err := o.apiClient.GetUsers()
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	for _, u := range users {
		if !o.perUserMatcher.MatchString(u.Username) {
			continue
		}
		if !o.collectedUsers[u.Username] {
			o.collectedUsers[u.Username] = true
			if err := o.addUserCharts(u); err != nil {
				o.Warning(err)
			}
		}
		mx[u.Username+"_bytes_received"] = u.BytesReceived
		mx[u.Username+"_bytes_sent"] = u.BytesSent
		mx[u.Username+"_connection_time"] = now - u.ConnectedSince
	}

	return nil
}

func (o *OpenVPN) addUserCharts(user client.User) error {
	cs := userCharts.Copy()

	for _, chart := range *cs {
		chart.ID = fmt.Sprintf(chart.ID, user.Username)
		chart.Fam = fmt.Sprintf(chart.Fam, user.Username)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, user.Username)
		}
		chart.MarkNotCreated()
	}

	return o.charts.Add(*cs...)
}
