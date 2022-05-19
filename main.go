package main

import (
	"bufio"
	"com.csion/tasks-worker/net"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if strings.Contains(os.Getenv("os"), "Windows") {
		initEnv()
		go doSync()      // 异步写入port并开启master节点监听
		net.HttpServer(0) // 开启监听
	} else {
		args := os.Args
		if args[len(args)-1] == "-d" {
			initEnv()
			go doSync()     // 异步写入port并开启master节点监听
			net.HttpServer(0) // 开启监听
		} else {
			startDaemon(args)
		}
	}
}

// 异步写入port并开启master节点监听
func doSync() {
	// 设置端口
	time.Sleep(1e6)
	wd, _ := os.Getwd()
	file, _ := os.Create(wd + "/taskCluster.port")
	defer file.Close()
	_, _ = file.Write([]byte(strconv.Itoa(net.Port)))

	go func() {
		// 循环监听
		client := &http.Client{
			Timeout: time.Second * 2,
		}
		for {
			time.Sleep(time.Second * 5)
			MasterProbe(3, client)
		}
	}()
}

// 对master的监听
func MasterProbe(i int, client *http.Client) {
	r, err := client.Get(fmt.Sprintf("http://%s/node/ping", os.Getenv("MNode")))
	if err != nil {
		if i == 1 {
			os.Exit(0)
		}
		MasterProbe(i - 1, client)
	} else {
		r.Body.Close()
	}
}

// 启动Daemon服务
func startDaemon(args []string) {
	var cmd *exec.Cmd
	if strings.Contains(os.Getenv("os"), "Windows"){
		var arg []string
		if len(args) > 1 {
			arg = args[1:]
		}
		arg = append(arg, "-d")
		cmd = exec.Command(args[0], arg...)
	} else {
		c := sourceEnv()
		for _, arg := range args {
			c.WriteString(arg)
			c.WriteString(" ")
		}
		c.WriteString("-d")

		cmd = exec.Command("/bin/bash", "-c", c.String())
	}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

// 宿主机环境
func sourceEnv() *strings.Builder {
	wd, _ := os.Getwd()
	s := []string{"/etc/profile", wd + "/.bash_profile"}
	sb := strings.Builder{}
	for _, v := range s {
		if _, err := os.Stat(v); err == nil {
			sb.WriteString("source ")
			sb.WriteString(v)
			sb.WriteString(" && ")
		}
	}
	return &sb
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
