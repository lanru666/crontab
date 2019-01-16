package worker

import (
	"fmt"
	"github.com/lanru666/crontab/common"
	"time"
)

//调度协程 负责任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent               // etcd事件队列
	jobPlanTable      map[string]*common.JobSchedulerPlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo   //任务执行表
	jobResultChan     chan *common.JobExecuteResult       //任务执行结果
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
	case common.JOB_EVENT_KILL: //强杀任务事件
	
		
	}
}

// 尝试执行任务，一个一个执行
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	//调度和执行是两件事情
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	//执行的任务可能运行很久,1分钟会调度60次，但是只能执行1次,防止并发
	
	//如果任务正在执行，跳过本次调度
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("尚未退出,跳过执行")
		return
	}
	// 构建执行状态信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo
	// 执行任务
	fmt.Println("执行任务", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)
	return
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
			scheduler.TryStartJob(jobPlan)
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
		jobResult     *common.JobExecuteResult
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
		case jobResult = <-scheduler.jobResultChan: //监听任务执行结果
			scheduler.handleJobResult(jobResult)
		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	//删除执行状态,从执行表删除
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)
	fmt.Println("任务执行完成", result.ExecuteInfo.Job.Name, result.OutPut, result.Err)
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	//启动调度协程
	go G_scheduler.schedulerLoop()
	return
}

// 回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
