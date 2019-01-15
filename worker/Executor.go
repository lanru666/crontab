package worker

import (
	"context"
	"github.com/lanru666/crontab/common"
	"os/exec"
)

// 任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd    *exec.Cmd
			err    error
			output []byte
			result *common.JobExecuteResult
		)
		// 执行shell命令
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		// 执行并捕获输出
		output, err = cmd.CombinedOutput()
		//任务执行完成后，把执行的结果返回给Scheduler,Scheduler会从ExecutingTable中删除掉执行记录
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			OutPut:      make([]byte, 0),
			
		}
	}()
}

//初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{
	
	}
	return
}
