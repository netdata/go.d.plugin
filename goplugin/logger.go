package goplugin

import "github.com/l2isbad/go.d.plugin/logger"

type Logger interface {
	Critical(...interface{})
	Error(...interface{})
	Warning(...interface{})
	Info(...interface{})
	Debug(...interface{})

	Criticalf(string, ...interface{})
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})

	SetLevel(logger.Severity)
	Level() logger.Severity
}
