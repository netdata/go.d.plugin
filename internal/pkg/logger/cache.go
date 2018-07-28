package logger

var cache = map[string]*Logger{}

func add(l *Logger) {
	cache[l.ModuleName()+l.JobName()] = l
}

func CacheGet(n namer) *Logger {
	return cache[n.ModuleName()+n.JobName()]
}
