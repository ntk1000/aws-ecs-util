package main

import (
	"os"
	"testing"
)

const (
	ExitNoEnv = "env not set"
	KEYID     = "AWS_ACCESS_KEY_ID"
	SECRET    = "AWS_SECRET_ACCESS_KEY"
	REGION    = "AWS_DEFAULT_REGION"
)

type conf struct {
	keyid, secret, region string
}

func confSetUp(t *testing.T) (c conf) {
	c.keyid = os.Getenv(KEYID)
	if c.keyid == "" {
		t.Errorf(ExitNoEnv)
	}
	c.secret = os.Getenv(SECRET)
	if c.secret == "" {
		t.Errorf(ExitNoEnv)
	}
	c.region = os.Getenv(REGION)
	if c.region == "" {
		t.Errorf(ExitNoEnv)
	}
	return
}

func TestCreateServiceClientViaEnv(t *testing.T) {

	c := confSetUp(t)
	svc, err := CreateServiceClientViaEnv()
	if svc == nil || err != nil {
		t.Errorf(ExitMsg, svc, "some service client")
		t.Errorf(ExitMsg, err, nil)
	}

	os.Setenv(KEYID, "")
	svc, err = CreateServiceClientViaEnv()
	if err == nil {
		t.Errorf(ExitMsg, err, "some error")
	}
	os.Setenv(KEYID, c.keyid)

	os.Setenv(SECRET, "")
	svc, err = CreateServiceClientViaEnv()
	if err == nil {
		t.Errorf(ExitMsg, err, "some error")
	}
	os.Setenv(SECRET, c.secret)

	os.Setenv(REGION, "")
	svc, err = CreateServiceClientViaEnv()
	if err == nil {
		t.Errorf(ExitMsg, err, "some error")
	}
	os.Setenv(REGION, c.region)
}

func TestListClusters(t *testing.T) {

	confSetUp(t)
	svc, err := CreateServiceClientViaEnv()
	if svc == nil || err != nil {
		t.Errorf(ExitMsg, svc, "some service client")
		t.Errorf(ExitMsg, err, nil)
	}

	_, err = ListClusters(svc)
	if err != nil {
		t.Errorf(ExitMsg, err, "some error")
	}

}

func TestListServices(t *testing.T) {

	confSetUp(t)
	svc, err := CreateServiceClientViaEnv()
	if svc == nil || err != nil {
		t.Errorf(ExitMsg, svc, "some service client")
		t.Errorf(ExitMsg, err, nil)
	}

	tmp := "x-stg-oneoff"
	c := &tmp
	_, err = ListServices(svc, c)
	if err != nil {
		t.Errorf(ExitMsg, err, "some error")
	}

}

func TestDescServices(t *testing.T) {

	confSetUp(t)
	svc, err := CreateServiceClientViaEnv()
	if svc == nil || err != nil {
		t.Errorf(ExitMsg, svc, "some service client")
		t.Errorf(ExitMsg, err, nil)
	}

	tmp := "x-stg-rp"
	c := &tmp
	s := []*string{c}
	_, err = DescServices(svc, c, s)
	if err != nil {
		t.Errorf(ExitMsg, err, "some error")
	}

}
