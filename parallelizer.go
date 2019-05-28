package main

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

func main() {
	//add verbose mode

	hiWhite := color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()

	flag.Parse()
	args := flag.Args()

	argsString := "[\"" + strings.Join(args, "\", \"") + "\"]"
	fmt.Fprintf(color.Output, hiWhite(fmt.Sprintf("Starting commands: %s\n---------------------\n", argsString)))

	procColors := []color.Attribute{
		color.FgRed,
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
		color.FgHiRed,
		color.FgHiGreen,
		color.FgHiYellow,
		color.FgHiBlue,
		color.FgHiMagenta,
		color.FgHiCyan}

	procColorsLength := len(procColors)

	cmdLength := len(args)
	wg.Add(cmdLength)

	processes := make(processDataMap, cmdLength)
	messageChannel := make(chan processOutput)

	for i, arg := range args {
		procColor := procColors[i%procColorsLength]
		process := newProcessData(i, arg, procColor)
		processes[i] = process

		go handleCommand(process, messageChannel)
	}

	maxNameLength := processes.getMaxNameLength()

	printDone := make(chan bool)
	go func(channel <-chan processOutput, done chan<- bool, maxLength int) {
		for message := range channel {
			procData := processes[message.procIndex]

			c := color.New(procData.color).SprintFunc()
			name := procData.name + strings.Repeat(" ", maxNameLength-len(procData.name))
			fmt.Fprintf(color.Output, "%s | %s\n", c(name), message.outputLine)
		}
		done <- true
	}(messageChannel, printDone, maxNameLength)

	wg.Wait()
	close(messageChannel)
	<-printDone

	fmt.Fprintf(color.Output, hiWhite("---------------------\nFinished all commands"))
}

func getExecCommand(procInfo *processData) *exec.Cmd {
	//check for os, use bash for linux

	cmd := exec.Command("cmd", "/c", procInfo.cmdWithArgs)

	return cmd
}

func handleCommand(procInfo *processData, messageChannel chan<- processOutput) {
	defer wg.Done()

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
