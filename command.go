package main

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func getExecCommand(procInfo *processData) *exec.Cmd {
	//check for os, use bash for linux

	interpreter := "cmd"
	interprenterArgs := []string{"/c"}
	changePath := "cd"
	changePathArgs := []string{"/d"}
	cmdLink := "&&"

	var cmdArray []string
	if len(procInfo.path) > 0 {
		cmdArray = append(cmdArray, changePath)
		cmdArray = append(cmdArray, changePathArgs...)
		cmdArray = append(cmdArray, procInfo.path)
		cmdArray = append(cmdArray, cmdLink)
	}
	cmdArray = append(cmdArray, procInfo.cmd)
	cmdArray = append(cmdArray, procInfo.args...)

	finalCmd := strings.Join(cmdArray, " ")

	args := append(interprenterArgs, finalCmd)
	cmd := exec.Command(interpreter, args...)

	return cmd
}

func handleCommand(procInfo *processData, messageChannel chan<- processOutput, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	cmd := getExecCommand(procInfo)

	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	cmd.Start()

	stdOutScanner := bufio.NewScanner(stdOut)
	stdOutScanner.Split(bufio.ScanLines)
	stdErrScanner := bufio.NewScanner(stdErr)
	stdErrScanner.Split(bufio.ScanLines)

	var outputWaitGroup sync.WaitGroup
	outputWaitGroup.Add(2)

	go func() {
		defer outputWaitGroup.Done()

		for stdOutScanner.Scan() {
			stdOutText := stdOutScanner.Text()
			messageChannel <- processOutput{procIndex: procInfo.index, outputLine: string(stdOutText)}
		}
	}()

	go func() {
		defer outputWaitGroup.Done()

		for stdErrScanner.Scan() {
			stdErrText := stdErrScanner.Text()
			messageChannel <- processOutput{procIndex: procInfo.index, outputLine: color.RedString(string(stdErrText))}
		}
	}()

	outputWaitGroup.Wait()
	cmd.Wait()
}
