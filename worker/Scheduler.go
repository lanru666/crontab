package worker

import (
	"fmt"
	"github.com/lanru666/crontab/common"
	"time"
)

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
	var (
		jobSchedulerPlan *common.JobSchedulerPlan
		jobExisted       bool
		err              error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulerPlan, err = common.BuildJobSchedulerPlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	//jobEvent.Job
	case common.JOB_EVENT_DELETE: //删除任务事件
		if jobSchedulerPlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// 计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulerPlan
		now      time.Time
		nearTime *time.Time
	)
	//如果任务表为空，随便睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}
	//当前时间
	now = time.Now()
	//1、遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO 尝试执行任务
			fmt.Println("执行任务", jobPlan.Job.Name)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //更新下次执行时间
		}
		// 统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
		//下次调度时间 (最近要执行的任务调度时间-当前时间)
		scheduleAfter = (*nearTime).Sub(now)
		return
	}
	//2、过期的任务立即执行
	
	//3、统计最近的要过期的任务的时间(N秒后过期)
	return
}

// 调度协程
func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	// 初始化一次(1秒)
	scheduleAfter = scheduler.TrySchedule()
	//调度的延时定时器
	scheduleTimer = time.NewTimer(scheduleAfter)
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //监听任务变化事件
			//对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: //最近的任务到期了
		
		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulerPlan),
	}
	//启动调度协程
	go G_scheduler.schedulerLoop()
	return
}
