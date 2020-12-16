#!/bin/bash
echo 0
ps -ef | grep "myserver" | awk '{print $2}' | xargs kill >/dev/null 2>&1
echo 1
nohup ./myserver &

