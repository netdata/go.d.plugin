// SPDX-License-Identifier: GPL-3.0-or-later

package whoisquery

import "fmt"

func (wq *WhoisQuery) collect() (map[string]int64, error) {
	remainingTime, err := wq.prov.remainingTime()
	if err != nil {
		return nil, fmt.Errorf("%v (source: %s)", err, wq.Source)
	}

	mx := make(map[string]int64)
	wq.collectExpiration(mx, remainingTime)
	return mx, nil
}

func (wq WhoisQuery) collectExpiration(mx map[string]int64, remainingTime float64) {
	mx["expiry"] = int64(remainingTime)
	mx["days_until_expiration_warning"] = wq.DaysUntilWarn
	mx["days_until_expiration_critical"] = wq.DaysUntilCrit

}
