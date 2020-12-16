# mkdir -p /data/log/server
cd ../cmd/myserver
go build -o ../../bin/myserver
cd ../../
nohup ./myserver &

