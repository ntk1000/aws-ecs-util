package main

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// CreateServiceClient returns ECS client via env
func CreateServiceClientViaEnv() (e *ecs.ECS, err error) {
	sess := session.Must(session.NewSession())
	creds := credentials.NewEnvCredentials()
	_, err = creds.Get()
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		return nil, errors.New("ENV AWS_DEFAULT_REGION doesn't defined")
	}
	e = ecs.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
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
