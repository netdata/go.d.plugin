package weblog

import "regexp"

type Parser interface {
	Parse(line string) (groupMap, bool)
}

type csvParser struct {
	pattern pattern
	reader  csvReader

	data groupMap
}

func (cp *csvParser) Parse(line string) (groupMap, bool) {
	// TODO: conversion to []byte should be fixed
	lines, err := cp.reader.readRecord([]byte(line))

	if err != nil {
		return nil, false
	}

	// NOTE: no index out of bound check
	for _, p := range cp.pattern {
		cp.data[p.name] = lines[p.index]
	}
	return cp.data, true
}

type regexpParser struct {
	re *regexp.Regexp

	data groupMap
}

func (rp *regexpParser) Parse(line string) (groupMap, bool) {
	lines := rp.re.FindStringSubmatch(line)

	if lines == nil {
		return nil, false
	}

	for i, v := range rp.re.SubexpNames()[1:] {
		rp.data[v] = lines[i+1]
	}

	return rp.data, true
}

type groupMap map[string]string

func (gm groupMap) has(key string) bool {
	_, ok := gm[key]
	return ok
}

func (gm groupMap) get(s string) string {
	return gm[s]
}
func (gm groupMap) lookup(key string) (string, bool) {
	v, ok := gm[key]
	return v, ok
}
