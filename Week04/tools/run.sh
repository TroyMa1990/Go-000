#GODEBUG=http2debug=2 go run localHttp.go -log_dir=log -alsologtostderr -c 15 -a :8001

mkdir -p /data/log/http/local

go build localHttp.go
nohup ./localHttp -log_dir=/data/log/http/local &
