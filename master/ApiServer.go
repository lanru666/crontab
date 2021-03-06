package master

import (
	"encoding/json"
	"fmt"
	"github.com/lanru666/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	//单例对象 如果方法要被别的包调用，首字母大写
	G_apiServer *ApiServer
)

// 保存任务接口
// POST job = {"name":"job1","command":"echo hello","cronExpr":"* * * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("handleJobSave")
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//1、解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//2、取表单的job字段
	postJob = req.PostForm.Get("job")
	//3、反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//4、保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 5、返回正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6、返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 强制杀死任务
// POST /job/kill  name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("handleJobKill")
	var (
		err   error
		name  string
		bytes []byte
	)
	//1、POST:a=1&b=2&c=3
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//2、取杀死的任务名
	name = req.PostForm.Get("name")
	//3、去杀死任务
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	//正常应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6、返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 删除任务接口
// POST /job/delete  name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("handleJobDelete")
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	//1、POST:a=1&b=2&c=3
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//2、取删除的任务名
	name = req.PostForm.Get("name")
	//3、去删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	//正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6、返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 列举所有crontab任务
// GET /job/list
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("handleJobList")
	var (
		jobList []*common.Job
		err     error
		bytes   []byte
	)
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}
	jobList = jobList
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6、返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpSerer     *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler //静态文件http回调
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	//index.html
	//静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))
	
	//启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	// 创建一个HTTP服务
	httpSerer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpSerer,
	}
	//启动了服务端
	go httpSerer.Serve(listener)
	return
}
