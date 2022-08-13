package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/byterocket/c4udit/analyzer"
)

func main() {
	flag.Parse()
	if *help {
		printHelpAndExit()
	}

	// Expect at least one user argument.
	if len(flag.Args()) == 0 && !*toc{
		printHelpAndExit()
	}

	// Run analyzer.
	report, err := analyzer.Run(
		analyzer.AllIssues(),
		flag.Args(),
	)
	if err != nil {
		printErrorAndExit(err)
	}

	if *saveToFile {
		// Save report in markdown format to file.
		err = ioutil.WriteFile(
			"c4udit-report.md",
			[]byte(report.Markdown(false)),
			0777,
		)
		if err != nil {
			printErrorAndExit(err)
		}
	} else if *toc {
		// Save report in markdown format to file.
		// fmt.Println(flag.Args()[0])
		// str, err1 := analyzer.ToC_Convertor(flag.Args()[0])
		str, err1 := analyzer.ToC_Convertor("c4udit-report.md")
		if err1 != nil {
			printErrorAndExit(err1)
		}
		err = ioutil.WriteFile(
			"c4udit-report-toc.md",
			[]byte(str),
			0777,
		)
		if err != nil {
			printErrorAndExit(err)
		}

	} else {
		// Print report to stdout.
		fmt.Println(report.String())

	}

}

// Flags
var (
	help       = flag.Bool("h", false, "Print help text.")
	saveToFile = flag.Bool("s", false, "Save report as file.")
	toc        = flag.Bool("t", false, "Save Report as file with Toc")
)

const helpText = `c4udit is a static analyzer for solidity contracts based on regexs.

It is capable of finding low risk issues and gas optimizations documented in
the c4-common-issues[1] repository.

Note that c4udit has a high rate of false positives. Check the results carefully!


Usage:
	c4udit [flags] [files...]

Flags:
	-h    Print help text.
	-s    Save report as file.
	-t    Save report as file with Toc ex: ./c4udit -t

`

func printHelpAndExit() {
	fmt.Print(helpText)
	os.Exit(0)
}

func printErrorAndExit(err error) {
	fmt.Println("c4checker Error:")
	fmt.Print(err.Error())
	os.Exit(1)
}
