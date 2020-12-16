package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	//"github.com/golang/glog"
)

var (
	listenPort string
)

func init() {
	flag.StringVar(&listenPort, "l", "8880", "Http Server Listen Port, Default: 8880")
}

func HttpEchoHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	response := map[string]interface{}{
		"ret":     0,
		"errcode": 0,
		"msg":     "Success",
	}

	t, _ := json.Marshal(response)
	writer.Write(t)
}

func main() {
	http.HandleFunc("/echo", HttpEchoHandler)
	http.ListenAndServe(":"+listenPort, nil)
}
