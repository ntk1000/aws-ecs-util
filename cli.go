package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type CLI struct {
	outStream, errStream io.Writer
	Command              string
	WithAll              bool
	WithError            bool
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

// supported commands, consts
const (
	TaskCommand   = "show-tasks"
	EventsCommand = "show-events"
)

// supported flags
var (
	Command     string
	WithAll     bool
	WithError   bool
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

	var out []CustomEcsService
	// run command
	switch c.Command {
	case TaskCommand:
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
							out = append(out, ConvertService(*service))
						}
					} else {
						out = append(out, ConvertService(*service))
					}
				}
			}
		}
		// TODO else
	case EventsCommand:
		io.WriteString(c.outStream, "not yet")
	}

	if c.WithSlack {
		s := &Slack{}
		err := s.Init()
		if err != nil {
			io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
			return ExitCodeSlackError
		}
		//s.Setup(out)
		s.SetupAttachments(CreateSlackAttachments(out))
		err = s.Post()
		if err != nil {
			io.WriteString(c.errStream, fmt.Sprintf("%v\n", err))
			return ExitCodeSlackError
		}
	} else {
		b, _ := json.Marshal(out)
		io.WriteString(c.outStream, string(b))
	}

	return ExitCodeOK
}

func CreateSlackAttachments(ces []CustomEcsService) (a SlackAttachment) {
	for _, c := range ces {
		cf := &SlackField{
			Title: "cluster",
			Value: c.Cluster,
			Short: false,
		}
		sf := &SlackField{
			Title: "service",
			Value: c.Service,
			Short: false,
		}
		tf := &SlackField{
			Title: "taskdef",
			Value: c.TaskDefinition,
			Short: false,
		}
		df := &SlackField{
			Title: "desired",
			Value: strconv.FormatInt(c.Desired, 10),
			Short: false,
		}
		pf := &SlackField{
			Title: "pending",
			Value: strconv.FormatInt(c.Pending, 10),
			Short: false,
		}
		rf := &SlackField{
			Title: "running",
			Value: strconv.FormatInt(c.Running, 10),
			Short: false,
		}

		a.Fields = append(a.Fields, *cf, *sf, *tf, *df, *pf, *rf)
	}
	return
}

func ConvertService(s ecs.Service) CustomEcsService {
	var clusterarn = strings.Split(*s.ClusterArn, "/")
	var taskdef = strings.Split(*s.TaskDefinition, "/")
	return CustomEcsService{
		Cluster:        clusterarn[len(clusterarn)-1],
		Service:        *s.ServiceName,
		TaskDefinition: taskdef[len(taskdef)-1],
		Desired:        *s.DesiredCount,
		Pending:        *s.PendingCount,
		Running:        *s.RunningCount,
	}
}

type CustomEcsService struct {
	Cluster        string `json:"cluster"`
	Service        string `json:"service"`
	TaskDefinition string `json:"taskdef"`
	Desired        int64  `json:"desired"`
	Pending        int64  `json:"pending"`
	Running        int64  `json:"running"`
}
