package modules

import (
	"path"
	"runtime"
)

type Creator func() Module

type CreatorsMap map[string]Creator

var Registry = make(CreatorsMap)

func Add(c Creator) {
	if name := getFileName(2); name != "" {
		Registry[name] = c
	}
}

func getFileName(skip int) (name string) {
	_, n, _, _ := runtime.Caller(skip)
	if n != "" {
		_, name = path.Split(n[:len(n)-3])
	}
	return
}
