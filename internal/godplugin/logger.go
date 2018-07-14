package godplugin

import "github.com/l2isbad/go.d.plugin/internal/pkg/logger"

type n struct{}

func (*n) GetModuleName() string {
	return "plugin"
}

func (*n) GetJobName() string {
	return "main"
}

var log = logger.New(&n{})
