package logger

var cache = map[string]*logger{}

func add(l *logger) {
	cache[l.ModuleName()+l.JobName()] = l
}

func CacheGet(n namer) *logger {
	return cache[n.ModuleName()+n.JobName()]
}
