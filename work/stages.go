package work

import (
	"com.csion/tasks-worker/execShell"
	"os"
	"strings"
)

// git交互下代码
func Git(url string, branch string, workDir string, ch *chan []byte){
	execShell.ExecShell("git init", workDir, ch)
	execShell.ExecShell("git remote add origin " + url, workDir, ch)
	execShell.ExecShell("git fetch origin", workDir, ch)
	execShell.ExecShell("git checkout -b " + branch + " origin/" + branch, workDir, ch)
}

// 执行脚本
func ExecScript(scripts string, scriptDir string, workDir string, ch *chan []byte){
	filePath := execShell.CreateTempShell(scriptDir, scripts)
	if strings.Contains(os.Getenv("os"), "Windows"){
		execShell.ExecShell(filePath, workDir, ch)
	} else {
		execShell.ExecShell("sh " + filePath, workDir, ch)
	}
	// execShell.DelFile(filePath)
}

// http调用
func HttpInvoke(url string, param string, t string){

}
