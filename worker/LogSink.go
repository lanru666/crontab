package worker

import (
	"github.com/lanru666/crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

//Mongodb存储日志
type LogSink struct {
	client        *mongo.Client
	logCollection *mongo.Collection
	logChan       chan *common.JobLog
}

var (
	//单例
	G_logSink *LogSink
)

func InitLogSink(err error) {
	var (
		client *mongo.Client
		option *options.ClientOptions
	)
	option = options.Client().SetConnectTimeout(
		time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond)
	if client, err = mongo.NewClientWithOptions(G_config.MongodbUri, option); err != nil {
		return
	}
	//选择db和collection
	G_logSink = &LogSink{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
		logChan:       make(chan *common.JobLog, 1000),
	}
	//启动mongodbc处理协程
}

//日志存储协程
func (logSink *LogSink) writeLoop() {

}
