#!/usr/bin/env bash
#export totalEvents=10000

export logLevel="debug"

#export concurrencyLevel=50

#export numberOfUsers=100

( cd $(dirname $0)
time java -server -Xmx1G -jar ./follower-maze-2.0.jar )
