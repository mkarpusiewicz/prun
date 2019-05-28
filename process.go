package main

import (
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type processData struct {
	path  string
	cmd   string
	args  []string
	name  string
	index int
	color color.Attribute
}

type processOutput struct {
	procIndex  int
	outputLine string
}

func newProcessData(index int, procCmd string, procColor color.Attribute) *processData {
	cmdSlice := strings.Split(procCmd, " ")
	cmdString := cmdSlice[0]
	cmdArgs := cmdSlice[1:]

	path, cmd := filepath.Split(cmdString)

	return &processData{
		index: index,
		path:  path,
		cmd:   cmd,
		args:  cmdArgs,
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
