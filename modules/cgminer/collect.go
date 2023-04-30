package cgminer

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (c *Cgminer) Init() bool {
	if c.Exec == "" {
		c.Exec = "cgminer"
	}
	return true
}

func (c *Cgminer) Check() bool {
	_, err := exec.LookPath(c.Exec)
	return err == nil
}

func (c *Cgminer) Charts() *module.Charts {
	return &module.Charts{
		{
			ID:    "hash_rate",
			Title: "Hash Rate",
			Units: "hash/s",
			Fam:   "hash rate",
			Ctx:   "cgminer.hash_rate",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "hash_rate_gpu0", Name: "GPU 0", Algo: module.Stacked},
				{ID: "hash_rate_gpu1", Name: "GPU 1", Algo: module.Stacked},
				{ID: "hash_rate_gpu2", Name: "GPU 2", Algo: module.Stacked},
				// Add more GPUs as needed
			},
		},
		{
			ID:    "accepted",
			Title: "Accepted Shares",
			Units: "shares",
			Fam:   "accepted shares",
			Ctx:   "cgminer.accepted",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "accepted_gpu0", Name: "GPU 0", Algo: module.Stacked},
				{ID: "accepted_gpu1", Name: "GPU 1", Algo: module.Stacked},
				{ID: "accepted_gpu2", Name: "GPU 2", Algo: module.Stacked},
				// Add more GPUs as needed
			},
		},
		{
			ID:    "rejected",
			Title: "Rejected Shares",
			Units: "shares",
			Fam:   "rejected shares",
			Ctx:   "cgminer.rejected",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "rejected_gpu0", Name: "GPU 0", Algo: module.Stacked},
				{ID: "rejected_gpu1", Name: "GPU 1", Algo: module.Stacked},
				{ID: "rejected_gpu2", Name: "GPU 2", Algo: module.Stacked},
		},
		{
			ID:    "temperature",
			Title: "Temperature",
			Units: "celsius",
			Fam:   "temperature",
			Ctx:   "cgminer.temperature",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "temperature_gpu0", Name: "GPU 0", Algo: module.Stacked},
				{ID: "temperature_gpu1", Name: "GPU 1", Algo: module.Stacked},
				{ID: "temperature_gpu2", Name: "GPU 2", Algo: module.Stacked},
				// Add more GPUs as needed
			},
		},
		{
			ID:    "fan_speed",
			Title: "Fan Speed",
			Units: "percent",
			Fam:   "fan speed",
			Ctx:   "cgminer.fan_speed",
			Type:  module.Stacked,
			Dims: module.Dims{
				{ID: "fan_speed_gpu0", Name: "GPU 0", Algo: module.Stacked},
				{ID: "fan_speed_gpu1", Name: "GPU 1", Algo: module.Stacked},
				{ID: "fan_speed_gpu2", Name: "GPU 2", Algo: module.Stacked},
				// Add more GPUs as needed
			},
		},
	}
}

func (c *Cgminer) Collect() map[string]int64 {
	collected := make(map[string]int64)

	scanner := bufio.NewScanner(c.readData())
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) < 2 {
			continue
		}

		metric := line[0]
		value, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			continue
		}

		switch metric {
		case "hash_rate_gpu0":
			collected[metric] = value
		case "hash_rate_gpu1":
			collected[metric] = value
		case "hash_rate_gpu2":
			collected[metric] = value
			// Add more GPUs as needed
		case "accepted_gpu0":
			collected[metric] = value
		case "accepted_gpu1":
			collected[metric] = value
		case "accepted_gpu2":
			collected[metric] = value
			// Add more GPUs as needed
		case "rejected_gpu0":
			collected[metric] = value
		case "rejected_gpu1":
      			collected[metric] = value
		case "rejected_gpu2":
			collected[metric] = value
			// Add more GPUs as needed
		case "temperature_gpu0":
			collected[metric] = value
		case "temperature_gpu1":
			collected[metric] = value
		case "temperature_gpu2":
			collected[metric] = value
			// Add more GPUs as needed
		case "fan_speed_gpu0":
			collected[metric] = value
		case "fan_speed_gpu1":
			collected[metric] = value
		case "fan_speed_gpu2":
			collected[metric] = value
			// Add more GPUs as needed
		}
	}

	return collected
}
