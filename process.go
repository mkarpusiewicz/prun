package main

import "github.com/fatih/color"

type processData struct {
	cmd   string
	name  string
	index int
	color color.Attribute
}

type processOutput struct {
	procIndex  int
	outputLine string
}

func newProcessData(index int, cmd string, procColor color.Attribute) *processData {
	// set only specific field value with field key
	return &processData{
		index: index,
		cmd:   cmd,
		name:  cmd,
		color: procColor}
}

type processDataMap map[int]*processData

func (processes *processDataMap) getMaxNameLength() int {
	maxLength := 0
	for _, procInfo := range *processes {
		nameLen := len(procInfo.name)
		if nameLen > maxLength {
			maxLength = nameLen
		}
	}
	return maxLength
}