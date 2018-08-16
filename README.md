# aws-ecs-util

[![Build Status](https://travis-ci.org/ntk1000/aws-ecs-util.svg?branch=master)](https://travis-ci.org/ntk1000/aws-ecs-util)
[![Coverage Status](https://coveralls.io/repos/github/ntk1000/aws-ecs-util/badge.svg?branch=master)](https://coveralls.io/github/ntk1000/aws-ecs-util?branch=master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fntk1000%2Faws-ecs-util.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fntk1000%2Faws-ecs-util?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/ntk1000/aws-ecs-util)](https://goreportcard.com/report/github.com/ntk1000/aws-ecs-util)

go utility tool for aws ecs (gofe)

## description

* show desired,pending,running ecs tasks in all/specific cluster/services
* detect desired > running ecs tasks in all/specific cluster/services
* show service events in all/specific cluster
* post to slack

## install

TODO

## usage

```
# show desired,pending,running ecs tasks in all clusters
gofe show-tasks -a

# detect desired > running ecs tasks in all clusters (detect errors)
gofe show-tasks -a -e

# filter cluster-name with sc option
gofe show-tasks -cn cluster-name

# filter service-name with ss option
gofe show-tasks -sn service-name

# show service events
gofe show-events

```

### ref

* https://deeeet.com/writing/2014/08/27/cli-reference/
* http://yapcasia.org/2014/talk/show/b49cc53a-027b-11e4-9357-07b16aeab6a4
* https://www.gnu.org/prep/standards/html_node/Option-Table.html#Option-Table

