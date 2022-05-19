package net

import (
	"com.csion/tasks-worker/vo"
	"com.csion/tasks-worker/work"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

const TASK = "/task"		// 任务发起
const STOP = "/stop"		// 终止任务
const PING = "/ping"		// 监听worker node
const REMOVE = "/remove"	// 删除worker node

var Port = 8911
// 如果用tcp的话，需要自己处理粘包问题，所以这里使用http协议进行交互
func HttpServer(p int) {
	if p > 1000 && p < 60000 {	// 65535
		Port = p
	}
	http.HandleFunc(TASK, taskHandle)
	http.HandleFunc(PING, pongHandle)
	http.HandleFunc(REMOVE, removeHandle)
	addListen()
}

func taskHandle(writer http.ResponseWriter,  request *http.Request) {
	defer request.Body.Close()

	if os.Getenv("Auth") != request.Header.Get("auth"){
		_ = json.NewEncoder(writer).Encode(vo.HandleResp{Code: 401, Msg: "Failed", Data: "Auth验证失败"})
	}

	var taskInfo vo.TaskVO
	if err := json.NewDecoder(request.Body).Decode(&taskInfo); err != nil {
		_ = json.NewEncoder(writer).Encode(vo.HandleResp{Code: 400, Msg: "Error", Data: "bad request：" + err.Error()})
	}
	go work.RunTask(&taskInfo)

	_ = json.NewEncoder(writer).Encode(vo.HandleResp{Code: 200, Msg: "Success", Data: "任务发起成功"})
}

func pongHandle(writer http.ResponseWriter,  request *http.Request) {
	defer request.Body.Close()
	_ = json.NewEncoder(writer).Encode(vo.HandleResp{Code: 200, Msg: "Success", Data: "pong"})
}

func removeHandle(writer http.ResponseWriter,  request *http.Request) {
	defer os.Exit(0)
	defer request.Body.Close()
	_ = json.NewEncoder(writer).Encode(vo.HandleResp{Code: 200, Msg: "Success", Data: "stop"})
}

// 端口筛选
func addListen() {
	if err := http.ListenAndServe("0.0.0.0:" + strconv.Itoa(Port), nil); err != nil {
		Port = Port + 1
		addListen()
	}
}
