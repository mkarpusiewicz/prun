package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

func main() {
	//add verbose mode

	supressedMessages := make(map[string]bool)
	supressedMessages["Terminate batch job (Y/N)?"] = true

	hiWhite := color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	output := make(chan string)
	outputDone := make(chan bool)

	go func() {
		sig := <-sigs
		output <- hiWhite(fmt.Sprintf("---------------------\nTerminating due to signal: %s\n", sig))
	}()

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

		// handle interrupt
		go handleCommand(process, messageChannel)
	}

	maxNameLength := processes.getMaxNameLength()

	go func(channel <-chan processOutput, done chan<- bool, maxLength int) {
		for message := range channel {
			test := strings.TrimRight(strings.TrimRight(strings.TrimRight(message.outputLine, "\r\n"), "\n"), " ")
			if _, ok := supressedMessages[test]; ok {
				continue
			}

			procData := processes[message.procIndex]

			c := color.New(procData.color).SprintFunc()
			name := procData.name + strings.Repeat(" ", maxNameLength-len(procData.name))
			output <- fmt.Sprintf("%s | %s\n", c(name), message.outputLine)
		}
	}(messageChannel, outputDone, maxNameLength)

	go func(output <-chan string, done chan<- bool) {
		for message := range output {
			fmt.Fprintf(color.Output, message)
		}
		done <- true
	}(output, outputDone)

	wg.Wait()
	close(messageChannel)

	output <- hiWhite("---------------------\nFinished all commands\n")
	close(output)
	<-outputDone
}

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

	//handle color output from command
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

	//send sigint or sigterm on closing

	outputWaitGroup.Wait()
	cmd.Wait()
}
