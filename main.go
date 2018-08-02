package main

import (
	"os"
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}

	if init := cli.Init(os.Args); init != ExitCodeOK {
		os.Exit(init)
	} else {
		os.Exit(cli.Run())
	}
}
