package ApiServer

import (
	"net"
	"net/http"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

// 保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	
}

//初始化服务
func initApiServer() (err error) {
	var (
		mux       *http.ServeMux
		listener  net.Listener
		httpSerer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	//启动TCP监听
	if listener, err = net.Listen("tcp", ":8070"); err != nil {
		return
	}
	// 创建一个HTTP服务
	httpSerer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}
	return
}
