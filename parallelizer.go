package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println("args:", args)

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

	processes := make([]*processData, cmdLength)

	for i, arg := range args {
		//fmt.Printf("%d: %s\n", i, arg)
		procColor := procColors[i%procColorsLength]
		process := newProcessData(i, arg, procColor)
		processes[i] = process

		go handleCommand(process)
	}

	wg.Wait()

	hiWhite := color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()

	fmt.Fprintf(color.Output, hiWhite("Finished all commands!"))
}

func handleCommand(procInfo *processData) {
	defer wg.Done()

	//procInfo.output <- fmt.Sprintf("Test output from: %s", procInfo.name)
	//	c := color.New(cmdColor).SprintFunc()
	//fmt.Fprintf(color.Output, "%v: %s\n", cmdColor, c(cmd))
}
