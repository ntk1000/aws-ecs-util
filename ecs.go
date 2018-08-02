package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func StdOutService(w io.Writer, s ecs.Service) {
	var clusterarn = strings.Split(*s.ClusterArn, "/")
	var taskdef = strings.Split(*s.TaskDefinition, "/")
	io.WriteString(w, fmt.Sprintf("%+v\t", clusterarn[len(clusterarn)-1]))
	io.WriteString(w, fmt.Sprintf("%+v\t", *s.ServiceName))
	io.WriteString(w, fmt.Sprintf("%+v\t", taskdef[len(taskdef)-1]))
	io.WriteString(w, fmt.Sprintf("%+v\t", *s.DesiredCount))
	io.WriteString(w, fmt.Sprintf("%+v\t", *s.PendingCount))
	io.WriteString(w, fmt.Sprintf("%+v\t\n", *s.RunningCount))
}

// CreateServiceClient returns ECS client via env
func CreateServiceClientViaEnv() (e *ecs.ECS, err error) {
	sess := session.Must(session.NewSession())
	creds := credentials.NewEnvCredentials()
	_, err = creds.Get()
	// exit if AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY env not defined
	//if err != nil {
	//	return nil, err
	//}
	e = ecs.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
	return
}

// ListClusters returns ECS clustersoutput
func ListClusters(svc *ecs.ECS) (list *ecs.ListClustersOutput, err error) {
	input := &ecs.ListClustersInput{}
	list, err = svc.ListClusters(input)
	return
}

// ListServices returns ECS serviceoutput
func ListServices(svc *ecs.ECS, cluster *string) (list *ecs.ListServicesOutput, err error) {
	input := &ecs.ListServicesInput{
		Cluster: cluster,
	}
	list, err = svc.ListServices(input)
	return
}

// DescServices returns ECS services
func DescServices(svc *ecs.ECS, cluster *string, service []*string) (desc *ecs.DescribeServicesOutput, err error) {

	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: service,
	}
	desc, err = svc.DescribeServices(input)
	// dont exit if cluster has no services
	return
}
