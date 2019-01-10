package master

import (
	"net"
	"net/http"
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
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	
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
	if listener, err = net.Listen("tcp", ":"+string(G_config.ApiPort)); err != nil {
		return
	}
	// 创建一个HTTP服务
	httpSerer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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
