# aws-ecs-util

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
gofe show-tasks -c cluster-name

# filter service-name with ss option
gofe show-tasks -s service-name

# show service events
gofe show-events

```

### ref

* https://deeeet.com/writing/2014/08/27/cli-reference/
* http://yapcasia.org/2014/talk/show/b49cc53a-027b-11e4-9357-07b16aeab6a4
* https://www.gnu.org/prep/standards/html_node/Option-Table.html#Option-Table

