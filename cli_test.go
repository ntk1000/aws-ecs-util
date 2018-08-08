package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

var cli *CLI

const (
	ExitMsg = "Exit actual: %v, expected: %v"
)

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

func setUp() {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli = &CLI{
		outStream: outStream,
		errStream: errStream,
	}
}

func TestCheckCommand(t *testing.T) {

	err := cli.CheckCommand(TaskCommand)
	if err != nil {
		t.Errorf(ExitMsg, err, nil)
	}

	err = cli.CheckCommand(EventsCommand)
	if err != nil {
		t.Errorf(ExitMsg, err, nil)
	}

	err = cli.CheckCommand("dummy")
	if err == nil {
		t.Errorf(ExitMsg, err, errors.New("sorry, unsupported command"))
	}

}

func TestInit_Args(t *testing.T) {
	args := strings.Split("gofe", " ")
	status := cli.Init(args)
	if status != ExitCodeArgsError {
		t.Errorf(ExitMsg, status, ExitCodeArgsError)
	}
}

func TestInit_Command(t *testing.T) {

	args := strings.Split("gofe "+TaskCommand, " ")
	status := cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.Command != TaskCommand {
		t.Errorf(ExitMsg, cli.Command, TaskCommand)
	}

	args = strings.Split("gofe "+EventsCommand, " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.Command != EventsCommand {
		t.Errorf(ExitMsg, cli.Command, EventsCommand)
	}

	args = strings.Split("gofe unsupported", " ")
	status = cli.Init(args)
	if status != ExitCodeCommandError {
		t.Errorf(ExitMsg, status, ExitCodeCommandError)
	}

}

// TODO test all pattern
func TestRun(t *testing.T) {
	args := strings.Split("gofe "+TaskCommand+" -a -e", " ")
	status := cli.Init(args)
	status = cli.Run()
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}

	args = strings.Split("gofe "+TaskCommand+" -a -wh", " ")
	status = cli.Init(args)
	status = cli.Run()
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}

	args = strings.Split("gofe "+TaskCommand+" -a", " ")
	status = cli.Init(args)
	status = cli.Run()
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}

	args = strings.Split("gofe "+TaskCommand+" -a -s", " ")
	status = cli.Init(args)
	status = cli.Run()
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}

	args = strings.Split("gofe "+EventsCommand, " ")
	status = cli.Init(args)
	status = cli.Run()
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}

}

func TestInit_Flags(t *testing.T) {

	args := strings.Split("gofe "+TaskCommand+" -a", " ")
	status := cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if !cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, true)
	}
	if cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, false)
	}
	if cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, false)
	}
	if cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, false)
	}

	args = strings.Split("gofe "+TaskCommand+" -e", " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, false)
	}
	if !cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, true)
	}
	if cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, false)
	}
	if cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, false)
	}

	args = strings.Split("gofe "+TaskCommand+" -wh", " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, false)
	}
	if cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, false)
	}
	if !cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, true)
	}
	if cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, false)
	}

	args = strings.Split("gofe "+TaskCommand+" -s", " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, false)
	}
	if cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, false)
	}
	if cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, false)
	}
	if !cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, true)
	}

	args = strings.Split("gofe "+TaskCommand+" -cn testcluster", " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, false)
	}
	if cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, false)
	}
	if cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, false)
	}
	if cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, false)
	}

	if cli.ClusterName != "testcluster" {
		t.Errorf(ExitMsg, cli.ClusterName, "testcluster")
	}

	args = strings.Split("gofe "+TaskCommand+" -sn testservice", " ")
	status = cli.Init(args)
	if status != ExitCodeOK {
		t.Errorf(ExitMsg, status, ExitCodeOK)
	}
	if cli.WithAll {
		t.Errorf(ExitMsg, cli.WithAll, false)
	}
	if cli.WithError {
		t.Errorf(ExitMsg, cli.WithError, false)
	}
	if cli.WithHeader {
		t.Errorf(ExitMsg, cli.WithHeader, false)
	}
	if cli.WithSlack {
		t.Errorf(ExitMsg, cli.WithSlack, false)
	}
	if cli.ServiceName != "testservice" {
		t.Errorf(ExitMsg, cli.ServiceName, "testservice")
	}
	//cmd := "gofe " + TaskCommand + " --afadf"
	//args := strings.Split(cmd, " ")
	//fmt.Printf("%v", args)
	//status := cli.Init(args)
	//fmt.Printf("status %v", status)
	//if status != 2 {
	//	t.Errorf(ExitMsg, status, ExitCodeParseFlagError)
	//}

}
