package worker

import "github.com/lanru666/crontab/common"

//调度协程 负责任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent // etcd事件队列
}

var (
	G_scheduler *Scheduler
)
//调度协程
func (scheduler *Scheduler) schedulerLoop() {

}

//初始化调度器
func initScheduler() {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
	}
	//启动调度协程
	go G_scheduler.schedulerLoop()
	
}
