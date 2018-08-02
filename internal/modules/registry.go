package modules

import (
	"path"
	"runtime"
)

type (
	Creator  func() Module
	Creators map[string]Creator
)

var Registry = Creators{}

func (c Creator) MakeModule() Module {
	return c()
}

func (c *Creators) Destroy() {
	for k := range *c {
		(*c)[k] = nil
	}
	*c = nil
}

func Add(c Creator) {
	name := getFileName(2)

	if name != "" {
		Registry[name] = c
	}
}

func getFileName(skip int) string {
	var name string
	_, n, _, _ := runtime.Caller(skip)

	if n != "" {
		_, name = path.Split(n[:len(n)-3])
	}
	return name
}
