ps -ef | grep localHttp | awk '{print $2}' | xargs kill >/dev/null 2>&1
