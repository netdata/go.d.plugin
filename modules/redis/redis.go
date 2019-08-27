package redis

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: false, // @TODO change it back
		},
		Create: func() module.Module { return New() },
	}

	module.Register("redis", creator)
}

const (
	defaultHost = "localhost"
	defaultPort = 6379
)

// New returns a pointer to a Redis instance with default values.
func New() *Redis {
	return &Redis{
		charts: charts.Copy(),
		Config: Config{
			Host:     defaultHost,
			Port:     defaultPort,
			Password: "",
			DbNum:    0,
		},
	}
}

// Config for the Redis module
type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DbNum    int    `yaml:"dbnum"`
}

// Redis module
type Redis struct {
	module.Base
	charts      *Charts
	client      *redis.Client
	bgSaveTime  int64
	Config      `yaml:",inline"`
	UpdateEvery int `yaml:"update_every"`
}

// Cleanup closes the connection to the Redis server
func (r *Redis) Cleanup() {
	if err := r.client.Close(); err != nil {
		r.Errorf("cleanup: error on closing the redis client [%+v]: %v", r.Config, err)
	}
}

// Init checks the configuration parameters and pings the Redis server
func (r *Redis) Init() bool {
	if r.Host == "" {
		r.Error("Host is not set")
		return false
	}

	if r.Port == 0 {
		r.Error("Port is not set")
		return false
	}

	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + strconv.Itoa(r.Port),
		Password: r.Password,
		DB:       r.DbNum,
	})

	_, err := client.Ping().Result()
	if err != nil {
		r.Errorf("init: error on starting the redis client [%+v]: %v", r.Config, err)
		return false
	}

	r.client = client
	r.bgSaveTime = 0

	r.Debugf("using config %+v", r.Config)

	return true
}

// Check makes sure that at least one metric is being collected
func (r *Redis) Check() bool {
	return len(r.Collect()) > 0
}

// Charts returns the Redis charts
func (r *Redis) Charts() *Charts {
	return r.charts
}

// Collect returns a map of metrics
func (r *Redis) Collect() map[string]int64 {
	metrics := make(map[string]int64)

	res, err := r.client.Do("INFO").Result()
	if err != nil {
		r.Errorf("Could not get Redis INFO: %v", err)
		return nil
	}

	err = parseMetrics(res.(string), r.charts, metrics, &r.bgSaveTime, r.UpdateEvery)
	if err != nil {
		r.Errorf("Got error while parsing metrics: %v", err)
		r.Debugf("INFO: %+v", res)
	}

	r.Debugf("Metrics: %+v", metrics)

	return metrics
}

var autoParseDims = []string{
	"total_commands_processed", "instantaneous_ops_per_sec",
	"used_memory", "used_memory_lua",
	"total_net_input_bytes", "total_net_output_bytes",
	"evicted_keys",
	"total_connections_received", "rejected_connections",
	"connected_clients", "blocked_clients",
	"connected_slaves",
	"rdb_changes_since_last_save",
	"uptime_in_seconds",
}

func parseMetrics(
	info string,
	charts *Charts,
	metrics map[string]int64,
	bgSaveTime *int64,
	updateEvery int,
) error {
	metricIdx, valueIdx := -1, -1
	re := regexp.MustCompile("(?P<metric>[a-z0-9_]+):(?P<value>.*[^\r\n])")
	for idx, group := range re.SubexpNames() {
		switch group {
		case "metric":
			metricIdx = idx
		case "value":
			valueIdx = idx
		}
	}

	data := make(map[string]string)
	for _, match := range re.FindAllStringSubmatch(string(info), -1) {
		data[match[metricIdx]] = match[valueIdx]
	}

	for _, metric := range autoParseDims {
		v, err := fetchFromData(data, metric)
		if err != nil {
			return err
		}

		metrics[metric] = v
	}

	// hit_rate calculation
	var err error
	metrics["hit_rate"], err = fetchHitRate(data)
	if err != nil {
		return fmt.Errorf("could not fetch hit rate: %v", err)
	}

	// keys_redis calculation
	if err = fetchKeysPerDB(data, metrics, charts); err != nil {
		return fmt.Errorf("could not fetch keys per db: %v", err)
	}

	// rdb_bgsave_in_progress calculation
	if v, ok := data["rdb_bgsave_in_progress"]; ok && v != "0" {
		*bgSaveTime += int64(updateEvery)
	} else {
		*bgSaveTime = 0
	}
	metrics["rdb_bgsave_in_progress"] = *bgSaveTime

	// rdb_last_bgsave_status calculation
	metrics["rdb_last_bgsave_status"] = 0
	if v, ok := data["rdb_last_bgsave_status"]; ok && v == "ok" {
		metrics["rdb_last_bgsave_status"] = 1
	}

	return nil
}

func fetchFromData(data map[string]string, key string) (int64, error) {
	if v, ok := data[key]; ok {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("could not convert %q from %T(%+v) to int: %v", key, v, v, err)
		}

		return int64(parsed), nil
	}

	return 0, fmt.Errorf("could not fetch %q", key)
}

func fetchHitRate(data map[string]string) (int64, error) {
	keySpaceHits, err := fetchFromData(data, "keyspace_hits")
	if err != nil {
		return 0, err
	}
	keySpaceMisses, err := fetchFromData(data, "keyspace_misses")
	if err != nil {
		return 0, err
	}

	if keySpaceHits > 0 || keySpaceMisses > 0 {
		return (keySpaceHits * 100) / (keySpaceHits + keySpaceMisses), nil
	}

	return 0, nil
}

func fetchKeysPerDB(data map[string]string, metrics map[string]int64, charts *Charts) error {
	if !charts.Has(keysRedisChartId) {
		return fmt.Errorf("could not find chart %q", keysRedisChartId)
	}

	keysChart := charts.Get(keysRedisChartId)

	keyMatch := regexp.MustCompile("^db[0-9]+$")
	valueMatch := regexp.MustCompile("^keys=([0-9]+),.*")

	for key, value := range data {
		if !keyMatch.MatchString(key) {
			continue
		}

		if !valueMatch.MatchString(value) {
			return fmt.Errorf("unexpected keyspace format: %q", value)
		}

		if !keysChart.HasDim(key) {
			dim := module.Dim{ID: key, Name: key}
			if err := keysChart.AddDim(&dim); err != nil {
				return fmt.Errorf("could not add dim %q to keys chart: %v", key, err)
			}
		}

		keys, err := strconv.Atoi(valueMatch.FindStringSubmatch(value)[1])
		if err != nil {
			return fmt.Errorf("could not parse keys from keyspace %q: %v", value, err)
		}

		metrics[key] = int64(keys)
	}

	return nil
}
