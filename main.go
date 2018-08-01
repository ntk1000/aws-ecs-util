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
	TaskCommand   = "show-tasks"
	EventsCommand = "show-events"
)

// command and options
var (
	Command     string
	CommandFlag = flag.String("c", "", "command")
	AllFlag     = flag.Bool("a", false, "use this option for all clusters")
	ErrorFlag   = flag.Bool("e", false, "detect desired > running ecs tasks")
	ClusterFlag = flag.String("cn", "", "filter by cluster name")
	ServiceFlag = flag.String("sn", "", "filter by service name")
)

func main() {
	Init()
	Run()
}

// error helper function
func exitErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// Init checks and parses args
func Init() {
	// check args
	if len(os.Args) < 2 {
		exitErrorf("Usage: gofe -p <command> <options>")
	}
	flag.Parse()
	PrintFlags()
	Command = *CommandFlag
	err := CheckCommand(Command)
	if err != nil {
		exitErrorf("%v", err)
	}
}

// CheckCommand checks command string
func CheckCommand(command string) error {
	switch command {
	case TaskCommand:
		return nil
	case EventsCommand:
		return nil
	default:
		return errors.New("unsupported command")
	}
}

func PrintFlags() {
	fmt.Println(*CommandFlag)
	fmt.Println(*AllFlag)
	fmt.Println(*ErrorFlag)
	fmt.Println(*ClusterFlag)
	fmt.Println(*ServiceFlag)
}

func StdOutService(service ecs.Service) {
	fmt.Printf("cluster = %+v\n", *service.ClusterArn)
	fmt.Printf("service = %+v\n", *service.ServiceName)
	fmt.Printf("taskdefinition = %+v\n", *service.TaskDefinition)
	fmt.Printf("desired = %+v\n", *service.DesiredCount)
	fmt.Printf("pending = %+v\n", *service.PendingCount)
	fmt.Printf("running = %+v\n", *service.RunningCount)
}

// Run example
func Run() {
	// create client
	svc := CreateServiceClientViaEnv()

	var clusters *ecs.ListClustersOutput
	// run command
	if Command == TaskCommand {
		if *AllFlag {
			fmt.Printf("all")
			clusters = ListClusters(svc)

			for _, cluster := range clusters.ClusterArns {
				services := ListServices(svc, cluster)
				descs := DescServices(svc, cluster, services.ServiceArns)
				if *ErrorFlag {
					for _, service := range descs.Services {
						if *service.DesiredCount != *service.RunningCount {
							StdOutService(*service)
						}
					}
				} else {
					for _, service := range descs.Services {
						StdOutService(*service)
					}
				}
			}
		}

	} else if Command == EventsCommand {
		fmt.Printf(Command)
	}

	// list clusters
	lsc := ListClusters(svc)

	// describe services
	for _, v := range lsc.ClusterArns {
		lss := ListServices(svc, v)
		dss := DescServices(svc, v, lss.ServiceArns)
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
	_, err := creds.Get()
	// exit if AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY env not defined
	if err != nil {
		exitErrorf("%v", err)
	}
	return ecs.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
}

// ListClusters returns ECS clustersoutput
func ListClusters(svc *ecs.ECS) (list *ecs.ListClustersOutput) {
	input := &ecs.ListClustersInput{}
	list, err := svc.ListClusters(input)
	if err != nil {
		exitErrorf("%v", err)
	}
	return
}

// ListServices returns ECS serviceoutput
func ListServices(svc *ecs.ECS, cluster *string) (list *ecs.ListServicesOutput) {
	input := &ecs.ListServicesInput{
		Cluster: cluster,
	}
	list, err := svc.ListServices(input)
	if err != nil {
		exitErrorf("%v", err)
	}
	return
}

// DescServices returns ECS services
func DescServices(svc *ecs.ECS, cluster *string, service []*string) (desc *ecs.DescribeServicesOutput) {

	// debug
	//fmt.Printf("svc %v\n", svc)
	//fmt.Printf("cluster %v\n", *cluster)
	//for _, v := range service {
	//	fmt.Printf("service %v\n", *v)
	//}

	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: service,
	}
	desc, err := svc.DescribeServices(input)
	if err != nil {
		// dont exit if cluster has no services
		fmt.Errorf("%v", err)
	}
	return
}
