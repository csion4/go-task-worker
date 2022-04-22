package work

import (
	utils "com.csion/tasks-worker/uitls"
	"com.csion/tasks-worker/vo"
	"com.csion/tasks-worker/ws"
	"fmt"
	"os"
	"strconv"
)

const OverFlag = "Over!"

func RunTask(taskInfo *vo.TaskVO) {
	stages := taskInfo.Stages
	taskCode := taskInfo.TaskCode
	recordId := taskInfo.RecordId

	// 初始化工作目录
	taskHome := os.Getenv("TaskHome")
	if taskHome == "" {
		taskHome = "/data/taskHome/"
	}
	workSpace := taskHome + "workspace/"


	ch := make(chan string, 1024)
	go ws.WebSocketClient(&ch, taskCode, recordId)

	// 创建目录
	if err := utils.CreateDir(workSpace + taskCode, 0666); err != nil {
		ch <- "【ERROR】 Init Workspace Error:" + err.Error() + "\n"
		ch <- OverFlag
		return
	}
	if err := utils.CreateDir(workSpace + taskCode + "@script", 0666); err != nil {
		ch <- "【ERROR】 Init Workspace Error:" + err.Error() + "\n"
		ch <- OverFlag
		return
	}

	for _, stageMap := range stages {
		for stage, env := range stageMap {
			switch stage {
			case 1:
				ch <- "----【start stage clone git project】---- \n"
				Git(env["gitUrl"], env["branch"], workSpace + taskCode, &ch)
				break
			case 2:
				ch <- "----【start stage exec script】---- \n"
				ExecScript(env["script"], workSpace + taskCode + "@script" , workSpace + taskCode, &ch)
				break
			default:
				ch <- fmt.Sprintf("----【unknown stageId %s】----%s", strconv.Itoa(stage), "\n")
				break
			}
		}
	}
	ch <- OverFlag
}
