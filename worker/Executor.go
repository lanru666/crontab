package worker

import (
	"github.com/lanru666/crontab/common"
	"os/exec"
	"time"
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
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		// 任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			OutPut:      make([]byte, 0),
		}
		// 初始化分布式锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		// 记录任务开始时间
		result.StartTime = time.Now()
		err = jobLock.TryLock()
		defer jobLock.Unlock() //释放锁
		if err != nil { //上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// 上锁成功后,重置任务启动时间
			result.StartTime = time.Now()
			// 执行shell命令
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
			// 执行并捕获输出
			output, err = cmd.CombinedOutput()
			// 记录任务结束时间
			result.EndTime = time.Now()
			result.OutPut = output
			result.Err = err
		}
		//任务执行完成后，把执行的结果返回给Scheduler,Scheduler会从ExecutingTable中删除掉执行记录
		G_scheduler.PushJobResult(result)
	}()
}

//初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{
	
	}
	return
}
