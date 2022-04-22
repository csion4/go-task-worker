package main

import (
	"bufio"
	"com.csion/tasks-worker/net"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if strings.Contains(os.Getenv("os"), "Windows") {
		initEnv()
		fmt.Println("start worker...")
		net.HttpServer(0) // 开启监听
	} else {
		args := os.Args
		if args[len(args)-1] == "-d" {
			initEnv()
			go setPort()      // 异步写入port
			net.HttpServer(0) // 开启监听
		} else {
			startDaemon(args)
		}
	}
}

// 设置端口
func setPort() {
	time.Sleep(1e6)
	wd, _ := os.Getwd()
	file, _ := os.Create(wd + "/taskCluster.port")
	defer file.Close()
	_, _ = file.Write([]byte(strconv.Itoa(net.Port)))
}

// 启动Daemon服务
func startDaemon(args []string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	arg = append(arg, "-d")
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

// 解析配置文件
func initEnv() {
	wd, _ := os.Getwd()
	conf, err := os.Open(wd + "/taskCluster.conf")
	if err != nil {
		return
	}
	defer conf.Close()
	r := bufio.NewReader(conf)
	for ;; {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			return
		}
		s := string(line)
		split := strings.Split(s, "=")
		if len(split) > 1 {
			_ = os.Setenv(split[0], split[1])
		}
	}
}
