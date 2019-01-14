package worker

import "github.com/lanru666/crontab/common"

//调度协程 负责任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent               // etcd事件队列
	jobPlanTable map[string]*common.JobSchedulerPlan //任务调度计划表
}

var (
	G_scheduler *Scheduler
)
// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
	//jobEvent.Job
	case common.JOB_EVENT_DELETE: //删除任务事件
	
	}
}

// 调度协程
func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //监听任务变化事件
			//对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		}
	}
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func initScheduler() {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
	}
	//启动调度协程
	go G_scheduler.schedulerLoop()
	
}
