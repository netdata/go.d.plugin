// SPDX-License-Identifier: GPL-3.0-or-later

package logger

import "strings"

var globalSeverity = INFO

// Severity is a logging severity level
type Severity int

const (
	// CRITICAL severity level
	CRITICAL Severity = iota
	// ERROR severity level
	ERROR
	// WARNING severity level
	WARNING
	// INFO severity level
	INFO
	// DEBUG severity level
	DEBUG
)

// String returns human-readable string
func (s Severity) String() string {
	switch s {
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	}
	return "UNKNOWN"
}

// ShortString returns human-readable short string
func (s Severity) ShortString() string {
	switch s {
	case CRITICAL:
		return "CRIT"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	}
	return "UNKNOWN"
}

// SetSeverity sets global severity level
func SetSeverity(severity Severity) {
	globalSeverity = severity
}

func SetSeverityByName(severity string) {
	switch strings.ToUpper(severity) {
	case "CRIT", "CRITICAL":
		globalSeverity = CRITICAL
	case "ERROR":
		globalSeverity = ERROR
	case "WARN", "WARNING":
		globalSeverity = WARNING
	case "INFO":
		globalSeverity = INFO
	case "DEBUG":
		globalSeverity = DEBUG
	}
}

func IsDebug() bool {
	return globalSeverity == DEBUG
}
