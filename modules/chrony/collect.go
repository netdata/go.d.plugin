package chrony

import (
	"fmt"
	"time"
)

func (c *Chrony) collect() (map[string]int64, error) {
	if c.client == nil {
		client, err := c.newClient(c)
		if err != nil {
			return nil, err
		}
		c.client = client
	}

	mx := make(map[string]int64)

	if err := c.collectTracking(mx); err != nil {
		return nil, err
	}
	if err := c.collectActivity(mx); err != nil {
		return mx, err
	}

	return mx, nil
}

const (
	// https://github.com/mlichvar/chrony/blob/7daf34675a5a2487895c74d1578241ca91a4eb70/ntp.h#L70-L75
	leapStatusNormal         = 0
	leapStatusInsertSecond   = 1
	leapStatusDeleteSecond   = 2
	leapStatusUnsynchronised = 3
)

func (c *Chrony) collectTracking(mx map[string]int64) error {
	// https://github.com/mlichvar/chrony/blob/5b04f3ca902e5d10aa5948fb7587d30b43941049/client.c#L2129
	tp, err := c.client.Tracking()
	if err != nil {
		return fmt.Errorf("error on collecting tracking: %v", err)
	}

	mx["stratum"] = int64(tp.Stratum)
	mx["leap_status_normal"] = boolToInt(tp.LeapStatus == leapStatusNormal)
	mx["leap_status_insert_second"] = boolToInt(tp.LeapStatus == leapStatusInsertSecond)
	mx["leap_status_delete_second"] = boolToInt(tp.LeapStatus == leapStatusDeleteSecond)
	mx["leap_status_unsynchronised"] = boolToInt(tp.LeapStatus == leapStatusUnsynchronised)
	mx["root_delay"] = tp.RootDelay.Int64()
	mx["root_dispersion"] = tp.RootDispersion.Int64()
	mx["skew"] = tp.SkewPpm.Int64()
	mx["last_offset"] = tp.LastOffset.Int64()
	mx["rms_offset"] = tp.RmsOffset.Int64()
	mx["update_interval"] = tp.LastUpdateInterval.Int64()
	// handle chrony restarts
	if tp.RefTime.Time().Year() != 1970 {
		mx["ref_measurement_time"] = time.Now().Unix() - tp.RefTime.Time().Unix()
	}
	mx["residual_frequency"] = tp.ResidFreqPpm.Int64()
	// https://github.com/mlichvar/chrony/blob/5b04f3ca902e5d10aa5948fb7587d30b43941049/client.c#L1706
	mx["current_correction"] = abs(tp.CurrentCorrection.Int64())
	mx["frequency"] = abs(tp.FreqPpm.Int64())

	return nil
}

func (c *Chrony) collectActivity(mx map[string]int64) error {
	// https://github.com/mlichvar/chrony/blob/5b04f3ca902e5d10aa5948fb7587d30b43941049/client.c#L2791
	ap, err := c.client.Activity()
	if err != nil {
		return fmt.Errorf("error on collecting activity: %v", err)
	}

	mx["online_sources"] = int64(ap.Online)
	mx["offline_sources"] = int64(ap.Offline)
	mx["burst_online_sources"] = int64(ap.BurstOnline)
	mx["burst_offline_sources"] = int64(ap.BurstOffline)
	mx["unresolved_sources"] = int64(ap.Unresolved)

	return nil
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
