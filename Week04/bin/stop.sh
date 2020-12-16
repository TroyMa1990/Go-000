ps -ef | grep "myserver" | awk '{print $2}' | xargs kill >/dev/null 2>&1
