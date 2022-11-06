package nvme

import (
	"path/filepath"
	"time"
)

func (n *NVMe) collect() (map[string]int64, error) {
	if n.exec == nil {
		v, err := n.newNVMeCLI(n.Config)
		if err != nil {
			return nil, err
		}
		n.exec = v
	}

	if now := time.Now(); now.Sub(n.listDevicesTime) > n.listDevicesEvery {
		n.listDevicesTime = now
		devices, err := n.exec.list()
		if err != nil {
			return nil, err
		}

		seen := make(map[string]bool)
		for _, v := range devices.Devices {
			device := extractDeviceFromPath(v.DevicePath)
			seen[device] = true

			if !n.devicePaths[v.DevicePath] {
				n.devicePaths[v.DevicePath] = true
				n.addDeviceCharts(device)
			}
		}
		for path := range n.devicePaths {
			device := extractDeviceFromPath(path)
			if !seen[device] {
				delete(n.devicePaths, device)
				n.removeDeviceCharts(device)
			}
		}
	}

	mx := make(map[string]int64)

	for path := range n.devicePaths {
		stats, err := n.exec.smartLog(path)
		if err != nil {
			return nil, err
		}

		device := extractDeviceFromPath(path)

		mx["device_"+device+"_temperature"] = int64(float64(stats.Temperature) - 273.15) // Kelvin => Celsius
		mx["device_"+device+"_percentage_used"] = stats.PercentUsed
		mx["device_"+device+"_available_spare"] = stats.AvailSpare
		mx["device_"+device+"_data_units_read"] = stats.DataUnitsRead * 1000 * 512
		mx["device_"+device+"_data_units_written"] = stats.DataUnitsWritten * 1000 * 512
		mx["device_"+device+"_host_read_commands"] = stats.HostReadCommands
		mx["device_"+device+"_host_write_commands"] = stats.HostWriteCommands
		mx["device_"+device+"_power_cycles"] = stats.PowerCycles
		mx["device_"+device+"_power_on_time"] = stats.PowerOnHours * 3600 // hours => seconds
		mx["device_"+device+"_unsafe_shutdowns"] = stats.UnsafeShutdowns
		mx["device_"+device+"_media_errors"] = stats.MediaErrors
		mx["device_"+device+"_num_err_log_entries"] = stats.NumErrLogEntries
		mx["device_"+device+"_controller_busy_time"] = stats.ControllerBusyTime * 60 // minutes => seconds
		mx["device_"+device+"_warning_temp_time"] = stats.WarningTempTime * 60       // minutes => seconds
		mx["device_"+device+"_critical_comp_time"] = stats.CriticalCompTime * 60     // minutes => seconds
		mx["device_"+device+"_thm_temp1_trans_count"] = stats.ThmTemp1TransCount
		mx["device_"+device+"_thm_temp2_trans_count"] = stats.ThmTemp2TransCount
		mx["device_"+device+"_thm_temp1_total_time"] = stats.ThmTemp1TotalTime // seconds
		mx["device_"+device+"_thm_temp2_total_time"] = stats.ThmTemp2TotalTime // seconds

		mx["device_"+device+"_critical_warning_available_spare"] = boolToInt(stats.CriticalWarning&1 != 0)
		mx["device_"+device+"_critical_warning_temp_threshold"] = boolToInt(stats.CriticalWarning&(1<<1) != 0)
		mx["device_"+device+"_critical_warning_nvm_subsystem_reliability"] = boolToInt(stats.CriticalWarning&(1<<2) != 0)
		mx["device_"+device+"_critical_warning_read_only"] = boolToInt(stats.CriticalWarning&(1<<3) != 0)
		mx["device_"+device+"_critical_warning_volatile_mem_backup_failed"] = boolToInt(stats.CriticalWarning&(1<<4) != 0)
		mx["device_"+device+"_critical_warning_persistent_memory_read_only"] = boolToInt(stats.CriticalWarning&(1<<5) != 0)
	}

	return mx, nil
}

func extractDeviceFromPath(devicePath string) string {
	_, name := filepath.Split(devicePath)
	return name
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}
