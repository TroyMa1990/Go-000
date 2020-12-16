#!/bin/bash
ps -ef | grep '/data/log/http/local' | awk '{print $2}' | xargs kill >/dev/null 2>&1
cd ../tools
echo "Start localHttp"
nohup ./localHttp -log_dir=/data/log/http/local &

