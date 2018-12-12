package parser

import "github.com/netdata/go.d.plugin/modules/weblog/pattern"

func NewCSVParser(pattern pattern.CSVPattern) Parser {
	return &csvParser{
		pattern: pattern,
		reader: csvReader{
			comma: ' ',
		},
		data: make(GroupMap),
	}
}

type csvParser struct {
	pattern pattern.CSVPattern
	reader  csvReader

	data GroupMap
}

func (cp *csvParser) Parse(line string) (GroupMap, bool) {
	// TODO: conversion to []byte should be fixed
	lines, err := cp.reader.readRecord([]byte(line))

	if err != nil {
		return nil, false
	}

	// NOTE: no index out of bound check
	for _, p := range cp.pattern {
		cp.data[p.Name] = lines[p.Index]
	}
	return cp.data, true
}
