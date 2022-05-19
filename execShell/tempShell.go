package execShell

import (
	"com.csion/tasks-worker/uitls"
	"os"
	"runtime"
	"strings"
)

// 创建临时shell脚本
func CreateTempShell(scriptDir string, scripts string, ch *chan string) string {
	var tempFile string
	if strings.Contains(os.Getenv("os"), "Windows") {
		tempFile = "/temp.bat"
	} else {
		tempFile = "/temp.sh"
	}

	file, err := os.OpenFile(scriptDir + tempFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		*ch <- "【ERROR】:创建临时shell脚本异常 " + err.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	defer file.Close()

	_, err = file.WriteString(scripts)
	if err != nil {
		*ch <- "【ERROR】:写入临时shell脚本异常 " + err.Error() + "\n"
		*ch <- utils.FailedFlag
		runtime.Goexit()
	}
	return file.Name()
}
func DelFile(file string){
	_ = os.Remove(file)
}
