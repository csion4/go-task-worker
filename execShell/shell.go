package execShell

import (
	"bufio"
	"com.csion/tasks-worker/uitls"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//执行shell脚本
func ExecShell(cmd string, dir string, ch *chan string) {
	if strings.HasSuffix(cmd, ".sh") {
		*ch <- "【shell】\n"
	} else {
		*ch <- "【script】: " + cmd + " \n"
	}

	var command *exec.Cmd
	if strings.Contains(os.Getenv("os"), "Windows"){
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("/bin/bash", "-c", cmd)
	}
	command.Dir = dir

	errPipe, err1 := command.StderrPipe()
	if err1 != nil {
		*ch <- "【ERROR】:获取脚本执行结果异常: " + err1.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	defer errPipe.Close()

	pipe, err1 := command.StdoutPipe()
	if err1 != nil {
		*ch <- "【ERROR】:获取脚本执行结果异常: " + err1.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	defer pipe.Close()

	if err2 := command.Start(); err2 != nil {
		*ch <- "【ERROR】:脚本执行异常: " + err2.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}

	errOut, err1 := ioutil.ReadAll(errPipe)
	if err1 != nil {
		*ch <- "【ERROR】:脚本执行异常: " + err1.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	if len(errOut) > 0 {
		// *ch <- string(errOut) + "\n"
		*ch <- "【ERROR】:脚本执行异常: " + string(errOut) + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}

	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			*ch <- "【ERROR】:脚本执行异常: " + err.Error() + "\n"
			*ch <- utils.FailedFlag
			runtime.Goexit()
		}
		*ch <- string(line) + "\n"
	}

}
