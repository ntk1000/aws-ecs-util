package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {

	// create client
	svc := CreateServiceClientViaEnv()

	// list clusters
	lsc, _ := LsClusters(svc)

	// describe services
	for _, v := range lsc.ClusterArns {
		lss, _ := LsServices(svc, v)
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

func LsClusters(svc *ecs.ECS) (*ecs.ListClustersOutput, error) {
	input := &ecs.ListClustersInput{}
	return svc.ListClusters(input)
}

func LsServices(svc *ecs.ECS, cluster *string) (*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: cluster,
	}
	return svc.ListServices(input)
}

func DescServices(svc *ecs.ECS, cluster *string, service []*string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: service,
	}
	return svc.DescribeServices(input)
}

type ServiceStatus struct {
	ClusterArn     *string
	ServiceName    *string
	TaskDefinition *string
	DesiredCount   *int64
	PendingCount   *int64
	RunningCount   *int64
}
