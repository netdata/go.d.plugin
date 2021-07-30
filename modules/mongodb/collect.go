package mongo

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func addIfExists(serverStatus map[string]interface{}, key string, ms map[string]int64) {
	mMap := serverStatus
	keys := strings.Split(key, ".")
	for _, k := range keys {
		val, ok := mMap[k]
		if !ok {
			return
		}
		switch t := val.(type) {
		case map[string]interface{}:
			mMap = val.(map[string]interface{})
		case int64:
			if _, ok := mMap[toID(k)]; ok {
				ms[key] = t
			}
		case int32:
			if _, ok := mMap[toID(k)]; ok {
				ms[key] = int64(t)
			}
		default:
			return
		}
	}
}

func iterateServerStatus(ms map[string]int64, status map[string]interface{}) {
	for _, chart := range serverStatusCharts {
		for _, dim := range chart.Dims {
			addIfExists(status, dim.ID, ms)
		}
	}
}

func (m *Mongo) serverStatusCollect(ms map[string]int64) {
	var status map[string]interface{}
	command := bson.D{{Key: "serverStatus", Value: 1}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.Timeout)
	defer cancel()
	err := m.client.Database(m.Config.name).RunCommand(ctx, command).Decode(&status)
	if err != nil {
		m.Errorf("error get server status from mongo: %s", err)
		return
	}
	iterateServerStatus(ms, status)
}
