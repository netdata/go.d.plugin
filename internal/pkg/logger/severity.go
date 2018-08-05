package logger

var sevLevel = INFO

type Severity int

const (
	CRITICAL Severity = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

func (s Severity) String() string {
	switch s {
	case CRITICAL:
		return "CRIT "
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARN "
	case INFO:
		return "INFO "
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
