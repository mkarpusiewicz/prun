package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
		//fmt.Printf("%d: %s\n", i, arg)
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

func getExecCommand(procCmd string) *exec.Cmd {
	//always add cmd or bash, check for os

	// cmdSlice := strings.Split(procCmd, " ")
	// cmdString := cmdSlice[0]
	// cmdArgs := cmdSlice[1:]

	cmd := exec.Command("cmd", "/c", procCmd)

	return cmd
}

func handleCommand(procInfo *processData, messageChannel chan<- processOutput) {
	defer wg.Done()

	cmd := getExecCommand(procInfo.cmd)

	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	cmd.Start()

	cmdOutput, _ := ioutil.ReadAll(stdOut)
	cmdError, _ := ioutil.ReadAll(stdErr)
	//cmdOutput, err := cmd.CombinedOutput()

	//read per line
	//  scanner := bufio.NewScanner(stderr)
	// scanner.Split(bufio.ScanLines)
	// for scanner.Scan() {
	//     m := scanner.Text()
	//     fmt.Println(m)
	// }

	// error red
	//remove empty lines

	if len(cmdOutput) > 0 {
		messageChannel <- processOutput{procIndex: procInfo.index, outputLine: string(cmdOutput)}
	}

	if len(cmdError) > 0 {
		messageChannel <- processOutput{procIndex: procInfo.index, outputLine: string(cmdError)}
	}

	// if err != nil {
	// 	messageChannel <- processOutput{procIndex: procInfo.index, outputLine: string(err.Error())}
	// }

	cmd.Wait()

	//message := fmt.Sprintf("Test output from: %s", procInfo.name)
}
