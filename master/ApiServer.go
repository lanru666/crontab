package master

import (
	"encoding/json"
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
// POST job = {"name":"job1","command":"echo hello","confExpr":"* * * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
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
ERR:
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux       *http.ServeMux
		listener  net.Listener
		httpSerer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
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
