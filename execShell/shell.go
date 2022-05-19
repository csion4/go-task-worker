package execShell

import (
	"bufio"
	"com.csion/tasks-worker/uitls"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//执行shell脚本
func ExecShell(cmd string, dir string, ch *chan string) {
	*ch <- "【script】: " + cmd + " \n"

	var command *exec.Cmd
	if strings.Contains(os.Getenv("os"), "Windows"){
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("/bin/bash", "-c", cmd)
	}
	command.Dir = dir

	pipe, err1 := command.StdoutPipe()
	if err1 != nil {
		*ch <- "【ERROR】:获取脚本执行结果异常" + err1.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	defer pipe.Close()

	if err2 := command.Start(); err2 != nil {
		*ch <- "【ERROR】:脚本执行异常" + err2.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}

	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			*ch <- "【ERROR】:脚本执行异常" + err.Error() + "\n"
			*ch <- utils.FailedFlag
			runtime.Goexit()
		}
		*ch <- string(line) + "\n"
	}

}
