package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
)

const queryVersion = "SELECT VERSION()"

var reVersion = regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)

func (m *MySQL) collectVersion() (ver *semver.Version, isMariaDB bool, err error) {
	// https://mariadb.com/kb/en/version/
	m.Debugf("executing query: '%s'", queryVersion)
	var verStr string
	if err := m.db.QueryRow(queryVersion).Scan(&verStr); err != nil {
		return nil, false, err
	}

	m.Debugf("application version: %s", verStr)

	// version string is not always valid semver (ex.: 8.0.22-0ubuntu0.20.04.2)
	match := reVersion.FindStringSubmatch(verStr)
	if len(match) == 0 {
		return nil, false, fmt.Errorf("couldn't parse version string '%s'", verStr)
	}

	ver, err = semver.New(fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3]))
	if err != nil {
		return nil, false, fmt.Errorf("couldn't parse version string '%s': %v", verStr, err)
	}
	isMariaDB = strings.Contains(verStr, "MariaDB")

	return ver, isMariaDB, nil
}
