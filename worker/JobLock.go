package worker

import "go.etcd.io/etcd/clientv3"

//分布式锁(TXN事务锁)
type JobLock struct {
	//etcd客户端
	kv      clientv3.KV
	lease   clientv3.Lease
	jobName string //锁哪个任务
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
	
}
