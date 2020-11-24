package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
)

const queryVersion = "SELECT VERSION()"

var reVersionCore = regexp.MustCompile(`^\d+\.\d+\.\d+`)

func (m *MySQL) collectVersion() (ver *semver.Version, isMariaDB bool, err error) {
	// https://mariadb.com/kb/en/version/
	m.Debugf("executing query: '%s'", queryVersion)
	var fullVersion string
	if err := m.db.QueryRow(queryVersion).Scan(&fullVersion); err != nil {
		return nil, false, err
	}

	m.Debugf("application version: %s", fullVersion)

	// version string is not always valid semver (ex.: 8.0.22-0ubuntu0.20.04.2)
	versionCore := reVersionCore.FindString(fullVersion)
	if versionCore == "" {
		return nil, false, fmt.Errorf("couldn't parse version string '%s'", fullVersion)
	}

	ver, err = semver.New(versionCore)
	if err != nil {
		return nil, false, fmt.Errorf("couldn't parse version string '%s': %v", fullVersion, err)
	}
	isMariaDB = strings.Contains(fullVersion, "MariaDB")

	return ver, isMariaDB, nil
}
