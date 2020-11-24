package couchbase

import (
	"fmt"
)

func (cb *Couchbase) collect() (map[string]int64, error) {
	collected := make(map[string]int64)
	err := cb.collectBasicStats(collected)
	if err != nil {
		return nil, fmt.Errorf("error on creating a connection: %v", err)
	}

	return collected, nil
}

func (cb Couchbase) collectBasicStats(collected map[string]int64) error {

	if cb.conn == nil {
		conn, err := cb.client.connect(cb.Config.URL, cb.Config.Username, cb.Config.Password)
		if err != nil {
			return fmt.Errorf("error on creating a connection: %v", err)
		}
		cb.conn = conn
	}

	p, err := cb.conn.GetPool("default")
	if err != nil {
		fmt.Printf("Node: %#v\n\n", err)
	}

	for _, b := range p.BucketMap {
		bs := b.BasicStats
		quota_percent_used_ := fmt.Sprintf("quota_used_%s", b.Name)
		collected[quota_percent_used_] = int64(bs["quotaPercentUsed"].(float64))

		ops_per_sec := fmt.Sprintf("ops_%s", b.Name)
		collected[ops_per_sec] = int64(bs["opsPerSec"].(float64))

		disk_fetches := fmt.Sprintf("fetches_%s", b.Name)
		collected[disk_fetches] = int64(bs["diskFetches"].(float64))

		item_count := fmt.Sprintf("item_count_%s", b.Name)
		collected[item_count] = int64(bs["itemCount"].(float64))

		disk_used := fmt.Sprintf("disk_%s", b.Name)
		collected[disk_used] = int64(bs["diskUsed"].(float64))

		data_used := fmt.Sprintf("data_%s", b.Name)
		collected[data_used] = int64(bs["dataUsed"].(float64))

		mem_used := fmt.Sprintf("mem_%s", b.Name)
		collected[mem_used] = int64(bs["memUsed"].(float64))

		num_non_resident := fmt.Sprintf("num_non_resident_%s", b.Name)
		collected[num_non_resident] = int64(bs["vbActiveNumNonResident"].(float64))
	}

	return nil
}

func (cb Couchbase) collectBucketNames() ([]string, error) {
	if cb.conn == nil {
		conn, err := cb.client.connect(cb.Config.URL, cb.Config.Username, cb.Config.Password)
		if err != nil {
			return nil, fmt.Errorf("error on creating a connection: %v", err)
		}
		cb.conn = conn
	}

	p, err := cb.conn.GetPool("default")
	if err != nil {
		fmt.Printf("Node: %#v\n\n", err)
	}
	var bucketNames []string
	for _, b := range p.BucketMap {
		bucketNames = append(bucketNames, b.Name)
	}
	return bucketNames, nil
}
