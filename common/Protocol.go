package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name     string `json:"name"`     // 任务名
	Command  string `json:"command"`  // shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// 任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 //要调度的的任务信息
	Expr     *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time            //下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job      *Job      //任务信息
	PlanTime time.Time //理论上的调度时间
	RealTime time.Time //实际的调度时间
}

//HTTP接口应答
type Response struct {
	Errno int         `json:"errorno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 1、定义一个response
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data
	// 2、序列化json
	resp, err = json.Marshal(response)
	
	return
}

//反序列化job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// 从etcd的Key中提取任务名
// /cron/jobs/job10,抹掉/cron/jobs/
func ExtractJobName(jobKey string) (string) {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// /cron/killer/job10,抹掉/cron/killer/
func ExtractKillerName(killerKey string) (string) {
	return strings.TrimPrefix(killerKey, JOB_KILL_DIR)
}

//任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo //执行状态
	OutPut      []byte          //脚本输出
	Err         error           //脚本错误原因
	StartTime   time.Time       //启动时间
	EndTime     time.Time       //结束时间
}

type JobEvent struct {
	EventType int // SAVE DELETE
	Job       *Job
}

//任务变化事件有两种，1) 更新任务 2) 删除任务
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

//构造任务执行计划
func BuildJobSchedulerPlan(job *Job) (jobSchedulerPLan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	//解析JOB的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	//生成任务调度计划对象
	jobSchedulerPLan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

//构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulerPlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, //计划调度时间
		RealTime: time.Now(),               //真实调度时间
	}
	return
}
