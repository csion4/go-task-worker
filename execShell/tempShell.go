package execShell

import (
	"os"
	"strings"
)

// 创建临时shell脚本
func CreateTempShell(scriptDir string, scripts string) string {
	var tempFile string
	if strings.Contains(os.Getenv("os"), "Windows") {
		tempFile = "/temp.bat"
	} else {
		tempFile = "/temp.sh"
	}

	file, _ := os.Create(scriptDir + tempFile)
	defer file.Close()
	_, _ = file.WriteString(scripts)
	return file.Name()
}
func DelFile(file string){
	_ = os.Remove(file)
}
