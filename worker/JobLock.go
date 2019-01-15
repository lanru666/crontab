package worker

import (
	"context"
	"github.com/lanru666/crontab/common"
	"go.etcd.io/etcd/clientv3"
)

//分布式锁(TXN事务锁)
type JobLock struct {
	//etcd客户端
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string             //锁哪个任务
	cancelFunc context.CancelFunc //用于终止自动续租
}

// 初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

// 尝试上锁函数，乐观锁
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse //只读chan
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	// 1、创建租约(5秒)
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	//创建context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	//租约ID
	leaseId = leaseGrantResp.ID
	//2、 自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}
	//3、处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: // 自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()
	//4、创建事务txn
	txn = jobLock.kv.Txn(context.TODO())
	// 锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName
	//5、事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).Else(clientv3.OpGet(lockKey))
	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}
	//6、成功返回，失败释放租约
	if txnResp.Succeeded { //锁被占用
		goto FAIL
	}
FAIL:
	cancelFunc()                                  //取消自动续租
	jobLock.lease.Revoke(context.TODO(), leaseId) // 释放租约
	return
}

//释放锁
func (jobLock *JobLock) Unlock() {
	
}
