package nvme

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

type nvmeDeviceList struct {
	Devices []struct {
		DevicePath   string `json:"DevicePath"`
		UsedBytes    int64  `json:"UsedBytes"`
		PhysicalSize int64  `json:"PhysicalSize"`
		SectorSize   int64  `json:"SectorSize"`
	}
}

type nvmeDeviceSmartLog struct {
	CriticalWarning    int64 `json:"critical_warning"`
	Temperature        int64 `json:"temperature"`
	AvailSpare         int64 `json:"avail_spare"`
	SpareThresh        int64 `json:"spare_thresh"`
	PercentUsed        int64 `json:"percent_used"`
	DataUnitsRead      int64 `json:"data_units_read"`
	DataUnitsWritten   int64 `json:"data_units_written"`
	HostReadCommands   int64 `json:"host_read_commands"`
	HostWriteCommands  int64 `json:"host_write_commands"`
	ControllerBusyTime int64 `json:"controller_busy_time"`
	PowerCycles        int64 `json:"power_cycles"`
	PowerOnHours       int64 `json:"power_on_hours"`
	UnsafeShutdowns    int64 `json:"unsafe_shutdowns"`
	MediaErrors        int64 `json:"media_errors"`
	NumErrLogEntries   int64 `json:"num_err_log_entries"`
	WarningTempTime    int64 `json:"warning_temp_time"`
	CriticalCompTime   int64 `json:"critical_comp_time"`
	ThmTemp1TransCount int64 `json:"thm_temp1_trans_count"`
	ThmTemp2TransCount int64 `json:"thm_temp2_trans_count"`
	ThmTemp1TotalTime  int64 `json:"thm_temp1_total_time"`
	ThmTemp2TotalTime  int64 `json:"thm_temp2_total_time"`
}

type nvmeCLIExec struct {
	sudoPath string
	nvmePath string
	timeout  time.Duration
}

func (n *nvmeCLIExec) list() (*nvmeDeviceList, error) {
	data, err := n.execute("list", "--output-format=json")
	if err != nil {
		return nil, err
	}

	var v nvmeDeviceList
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (n *nvmeCLIExec) smartLog(devicePath string) (*nvmeDeviceSmartLog, error) {
	data, err := n.execute("smart-log", devicePath, "--output-format=json")
	if err != nil {
		return nil, err
	}

	var v nvmeDeviceSmartLog
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (n *nvmeCLIExec) execute(arg ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	if n.sudoPath != "" {
		args := append([]string{"-n", n.nvmePath}, arg...)
		return exec.CommandContext(ctx, n.sudoPath, args...).Output()
	}

	return exec.CommandContext(ctx, n.nvmePath, arg...).Output()
}
