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

### todolist

* minify docker image using [dockerslim](https://github.com/docker-slim/docker-slim)
* speed up with concurrency or other technique
* setup binary to github release

### ref

* cli design
	* https://deeeet.com/writing/2014/08/27/cli-reference/
	* http://yapcasia.org/2014/talk/show/b49cc53a-027b-11e4-9357-07b16aeab6a4
	* https://www.gnu.org/prep/standards/html_node/Option-Table.html#Option-Table
* performance tuning
	* https://qiita.com/naoina/items/d71ddfab31f4b29f6693#%E4%BE%8B%E3%81%88%E3%81%B0channel-%E3%82%92%E4%BD%BF%E3%82%8F%E3%81%AA%E3%81%84
	* http://blog.monochromegane.com/blog/2015/12/15/how-to-speed-up-the-platinum-searcher-v2/
	* https://mattn.kaoriya.net/software/lang/go/20180531104907.htm
	* https://mattn.kaoriya.net/software/lang/go/20161019124907.htm
	* http://kenzo0107.hatenablog.com/entry/2016/11/21/142248
	* https://go-tour-jp.appspot.com/concurrency/1

