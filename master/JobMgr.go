package master

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 任务管理器
type jobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	//单例
	G_jobMgr *jobMgr
)

func InitJobMgr(err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints:   []string{""},            //集群地址
		DialTimeout: 5000 * time.Millisecond, //连接超时
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	// 建立KV和lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	// 赋值单例
	G_jobMgr = &jobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}
