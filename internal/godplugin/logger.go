package godplugin

import "github.com/l2isbad/go.d.plugin/internal/pkg/logger"

type n struct{}

func (*n) ModuleName() string {
	return "plugin"
}

func (*n) JobName() string {
	return "main"
}

var log = logger.New(&n{})
