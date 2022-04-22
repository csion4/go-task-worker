package execShell

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

var OverFlag = []byte("Over!")
//执行shell脚本
func ExecShell(cmd string, dir string, ch *chan []byte) {
	 *ch <- []byte("【script】: " + cmd + " \n")

	 var command *exec.Cmd
	if strings.Contains(os.Getenv("os"), "Windows"){
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("/bin/sh", "-c", cmd)
	}
	// command := exec.Command("cmd", "/C", cmd)	window
	command.Dir = dir

	pipe, err1 := command.StdoutPipe()
	if err1 != nil {
		*ch <- []byte("【ERROR】:" + err1.Error() + "\n")
		*ch <- OverFlag
		return
	}
	defer pipe.Close()

	if err2 := command.Start(); err2 != nil {
		*ch <- []byte("【ERROR】:" + err2.Error() + "\n")
		*ch <- OverFlag
		return
	}

	reader := bufio.NewReader(pipe)
	for ;; {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		*ch <- append(line, ' ', '\n')
	}

}