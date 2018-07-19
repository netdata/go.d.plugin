package logger

var cache loggers

type loggers []*logger

func (l *loggers) get(n namer) *logger {
	for _, logger := range *l {
		if n.ModuleName() == logger.ModuleName() && n.JobName() == logger.JobName() {
			return logger
		}
	}
	return nil
}

func (l *loggers) add(n *logger) {
	*l = append(*l, n)
}

func CacheGet(n namer) *logger {
	return cache.get(n)
}
