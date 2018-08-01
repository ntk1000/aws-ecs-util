package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

const (
	Version       = "v 0.0.1"
	TaskCommand   = "show-tasks"
	EventsCommand = "show-events"
)

// command and options
var (
	Command     string
	AllFlag     = flag.Bool("a", false, "use this option for all clusters")
	ErrorFlag   = flag.Bool("e", false, "detect desired > running ecs tasks")
	ClusterFlag = flag.String("c", "", "filter by cluster name")
	ServiceFlag = flag.String("s", "", "filter by service name")
	VersionFlag = flag.String("v", "", Version)
)

func exitErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {

	Init()

	Run()

}

// Init checks and parses args
func Init() {
	// check args
	if len(os.Args) < 2 {
		exitErrorf("Usage: gofe <command> <options>")
	}

	Command = os.Args[1]
	err := CheckCommand(Command)
	if err != nil {
		exitErrorf("%v", err)
	}

	flag.Parse()
}

// CheckCommand checks command string
func CheckCommand(command string) error {
	switch command {
	case TaskCommand:
		return nil
	case EventsCommand:
		return nil
	case "-h":
		return nil
	case "-v":
		return nil
	default:
		return errors.New("unsupported command")
	}

}

// Run example
// TODO error handling
func Run() {
	// create client
	svc := CreateServiceClientViaEnv()

	// list clusters
	lsc, _ := ListClusters(svc)

	// describe services
	for _, v := range lsc.ClusterArns {
		lss, _ := ListServices(svc, v)
		dss, _ := DescServices(svc, v, lss.ServiceArns)
		for _, vv := range dss.Services {
			// stdout if desired <> running
			if *vv.DesiredCount != *vv.RunningCount {
				fmt.Printf("cluster = %+v\n", *vv.ClusterArn)
				fmt.Printf("service = %+v\n", *vv.ServiceName)
				fmt.Printf("taskdefinition = %+v\n", *vv.TaskDefinition)
				fmt.Printf("desired = %+v\n", *vv.DesiredCount)
				fmt.Printf("pending = %+v\n", *vv.PendingCount)
				fmt.Printf("running = %+v\n", *vv.RunningCount)
			}
		}
	}

}

// CreateServiceClient returns ECS client via env
func CreateServiceClientViaEnv() *ecs.ECS {
	sess := session.Must(session.NewSession())
	creds := credentials.NewEnvCredentials()
	creds.Get()
	return ecs.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
}

// ListClusters returns ECS clustersoutput
func ListClusters(svc *ecs.ECS) (*ecs.ListClustersOutput, error) {
	input := &ecs.ListClustersInput{}
	return svc.ListClusters(input)
}

// ListServices returns ECS serviceoutput
func ListServices(svc *ecs.ECS, cluster *string) (*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: cluster,
	}
	return svc.ListServices(input)
}

// DescServices returns ECS services
func DescServices(svc *ecs.ECS, cluster *string, service []*string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: service,
	}
	return svc.DescribeServices(input)
}

// ServiceStatus is minimum struct of ECS ServiceStatus
type ServiceStatus struct {
	ClusterArn     *string
	ServiceName    *string
	TaskDefinition *string
	DesiredCount   *int64
	PendingCount   *int64
	RunningCount   *int64
}
