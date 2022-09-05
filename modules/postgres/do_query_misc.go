// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import "strconv"

func (p *Postgres) doQuerySettingsMaxConnections() (int64, error) {
	q := querySettingsMaxConnections()

	var s string
	if err := p.doQueryRow(q, &s); err != nil {
		return 0, err
	}

	return strconv.ParseInt(s, 10, 64)
}

func (p *Postgres) doQueryServerVersion() (int, error) {
	q := queryServerVersion()

	var s string
	if err := p.doQueryRow(q, &s); err != nil {
		return 0, err
	}

	return strconv.Atoi(s)
}

func (p *Postgres) doQueryIsSuperUser() (bool, error) {
	q := queryIsSuperUser()

	var v bool
	if err := p.doQueryRow(q, &v); err != nil {
		return false, err
	}

	return v, nil
}
