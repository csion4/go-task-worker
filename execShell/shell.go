package execShell

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

const OverFlag = "Over!"
//执行shell脚本
func ExecShell(cmd string, dir string, ch *chan string) {
	*ch <- "【script】: " + cmd + " \n"

	var command *exec.Cmd
	if strings.Contains(os.Getenv("os"), "Windows"){
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("/bin/sh", "-c", cmd)
	}
	command.Dir = dir

	pipe, err1 := command.StdoutPipe()
	if err1 != nil {
		*ch <- "【ERROR】:获取脚本执行结果异常" + err1.Error() + "\n"
		*ch <- OverFlag
		panic(err1)
	}
	defer pipe.Close()

	if err2 := command.Start(); err2 != nil {
		*ch <- "【ERROR】:脚本执行异常" + err2.Error() + "\n"
		*ch <- OverFlag
		panic(err2)
	}

	reader := bufio.NewReader(pipe)
	for ;; {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			*ch <- "【ERROR】:脚本执行异常" + err.Error() + "\n"
			*ch <- OverFlag
			panic(err)
		}
		*ch <- string(line) + "\n"
	}

}
