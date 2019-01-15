package worker

import "github.com/lanru666/crontab/common"

// 任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {

}

//初始化执行器
func InitExecutor() (err error) {
	return
}
