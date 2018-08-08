package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type CLI struct {
	outStream, errStream io.Writer
	Command              string
	WithAll              bool
	WithError            bool
	WithHeader           bool
	ClusterName          string
	ServiceName          string
	WithSlack            bool
}

// exit codes
const (
	ExitCodeOK = iota
	ExitCodeArgsError
	ExitCodeCommandError
	ExitCodeParseFlagError
	ExitCodeConfigError
	ExitCodeAWSECSError
	ExitCodeSlackError
)

// supported commands
const (
	TaskCommand   = "show-tasks"
	EventsCommand = "show-events"
)

// supported flags
var (
	Command     string
	WithAll     bool
	WithError   bool
	WithHeader  bool
	ClusterName string
	ServiceName string
	WithSlack   bool
)

// Init checks and parses args
func (c *CLI) Init(args []string) int {

	// check args
	if len(args) < 2 {
		return ExitCodeArgsError
	}

	c.Command = args[1]
	if err := c.CheckCommand(c.Command); err != nil {
		io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
		return ExitCodeCommandError
	}

	// setup flags
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&c.WithAll, "a", false, "use this option for all clusters")
	flags.BoolVar(&c.WithError, "e", false, "detect desired > running ecs tasks, error status")
	flags.BoolVar(&c.WithHeader, "wh", false, "show headers")
	flags.StringVar(&c.ClusterName, "cn", "", "target cluster name")
	flags.StringVar(&c.ServiceName, "sn", "", "target service name")
	flags.BoolVar(&c.WithSlack, "s", false, "use this option to send result to slack")

	// skip args[1] for command
	if err := flags.Parse(args[2:]); err != nil {
		return ExitCodeParseFlagError
	}

	return ExitCodeOK
}

// CheckCommand checks command string
func (c *CLI) CheckCommand(command string) error {
	switch command {
	case TaskCommand:
		return nil
	case EventsCommand:
		return nil
	default:
		return errors.New("sorry, unsupported command")
	}
}

// Run
func (c *CLI) Run() int {
	// create client
	svc, err := CreateServiceClientViaEnv()
	if err != nil {
		io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
		return ExitCodeConfigError
	}

	var out string
	// run command
	switch c.Command {
	case TaskCommand:
		if c.WithHeader {
			io.WriteString(c.outStream,
				"cluster\tservice\ttaskdefinition\tdesired\tpending\trunning\n")
		}
		if c.WithAll {
			clusters, err := ListClusters(svc)
			if err != nil {
				io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
				return ExitCodeAWSECSError
			}
			for _, cluster := range clusters.ClusterArns {
				services, err := ListServices(svc, cluster)
				if err != nil {
					io.WriteString(c.errStream, fmt.Sprintf("%v\n", *cluster))
					io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
					return ExitCodeAWSECSError
				}
				descs, err2 := DescServices(svc, cluster, services.ServiceArns)
				if err2 != nil {
					io.WriteString(c.errStream, fmt.Sprintf("%v\n", *cluster))
					io.WriteString(c.errStream, fmt.Sprintf("%v\n", err2))
					//return ExitCodeAWSECSError
				}
				for _, service := range descs.Services {
					if c.WithError {
						if *service.DesiredCount != *service.RunningCount {
							out += StdOutService(c.outStream, *service)
						}
					} else {
						out += StdOutService(c.outStream, *service)
					}
				}
			}
		}
		// TODO else
	case EventsCommand:
		io.WriteString(c.outStream, "not yet")
	}

	if c.WithSlack {
		s := NewSlack("", "", out)
		err := s.Post()
		if err != nil {
			io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
			return ExitCodeSlackError
		}
	}

	return ExitCodeOK
}
