package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

func main() {
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

func handleCommand(procInfo *processData, messageChannel chan<- processOutput) {
	defer wg.Done()

	message := fmt.Sprintf("Test output from: %s", procInfo.name)
	messageChannel <- processOutput{procIndex: procInfo.index, outputLine: message}
}
