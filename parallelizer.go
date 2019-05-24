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

	// fmt.Println(procColorsLength)

	// yellow := color.New(procColors[0]).SprintFunc()
	// red := color.New(procColors[6]).SprintFunc()

	// fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

	wg.Add(len(args))

	for i, arg := range args {
		//fmt.Printf("%d: %s\n", i, arg)
		color := procColors[i%procColorsLength]
		go handleCommand(arg, color)
	}

	wg.Wait()

	hiWhite := color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()

	fmt.Fprintf(color.Output, hiWhite("Finished all commands!"))
}

func handleCommand(cmd string, cmdColor color.Attribute) {
	defer wg.Done()

	c := color.New(cmdColor).SprintFunc()
	fmt.Fprintf(color.Output, "%v: %s\n", cmdColor, c(cmd))
}
