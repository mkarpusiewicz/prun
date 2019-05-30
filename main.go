package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	//todo: add verbose mode
	settings := getSettings()
	hiWhite := color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()

	//output handling
	output := make(chan string)
	outputDone := make(chan bool)

	go func(output <-chan string, done chan<- bool) {
		for message := range output {
			fmt.Fprintf(color.Output, message)
		}
		done <- true
	}(output, outputDone)

	//flag parsing
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		output <- fmt.Sprintf("\nprun ver: %v, commit: %v, built at %v\n", color.CyanString(version), color.CyanString(commit), color.CyanString(date))
		output <- fmt.Sprintf("\n   usage: %s %s %s %s\n", hiWhite("prun"), color.GreenString("cmd1"), color.YellowString("cmd2"), color.BlueString("cmd3"))

		close(output)
		<-outputDone

		os.Exit(0)
	}

	output <- hiWhite(fmt.Sprintf("Starting commands: [\"%s\"]\n---------------------\n", strings.Join(args, "\", \"")))

	//signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(output chan<- string) {
		sig := <-sigs
		output <- hiWhite(fmt.Sprintf("---------------------\nTerminating due to signal: %s\n", sig))
	}(output)

	//process handling
	cmdLength := len(args)
	processes := make(processDataMap, cmdLength)

	var wg sync.WaitGroup
	wg.Add(cmdLength)

	messageChannel := make(chan processOutput)

	for i, arg := range args {
		procColor := settings.GetProcessColor(i)
		process := newProcessData(i, arg, procColor)
		processes[i] = process

		go handleCommand(process, messageChannel, &wg)
	}

	maxNameLength := processes.getMaxNameLength()

	go func(channel <-chan processOutput, output chan<- string, done chan<- bool, maxLength int) {
		for message := range channel {
			if settings.ShouldSuppressMessage(message.outputLine) {
				continue
			}

			procData := processes[message.procIndex]

			c := color.New(procData.color).SprintFunc()
			name := procData.name + strings.Repeat(" ", maxLength-len(procData.name))

			output <- fmt.Sprintf("%s | %s\n", c(name), message.outputLine)
		}
	}(messageChannel, output, outputDone, maxNameLength)

	wg.Wait()
	close(messageChannel)

	output <- hiWhite("---------------------\nFinished all commands\n")
	close(output)
	<-outputDone
}
