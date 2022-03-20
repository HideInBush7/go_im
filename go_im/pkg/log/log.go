package log

import (
	"fmt"
	"sync"

	"github.com/HideInBush7/go_im/pkg/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type logConf struct {
	// 为true开启控制台输出
	Debug bool `json:"debug"`
	// 日志等级
	LogLevel string `json:"loglevel"`
	// 日志存储文件
	LogFile string `json:"logfile"`
}

var once sync.Once

func Init() {
	once.Do(func() {
		var conf = logConf{}
		config.Register("log", &conf)
		fmt.Println(conf)
		// 日志等级
		level, err := logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			panic(err)
		}
		logrus.SetLevel(level)

		// 设置软链和日志分割
		writer, err := rotatelogs.New(
			conf.LogFile+"%Y-%m-%d"+".log",
			rotatelogs.WithLinkName("log.log"),
		)
		if err != nil {
			logrus.Panic(err)
		}

		// 打开控制台输出 通过设置hook
		if conf.Debug {
			hook := lfshook.NewHook(writer, &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				PrettyPrint:     false, //是否格式化json格式
			})
			logrus.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
				DisableColors:   false,
			})
			logrus.AddHook(hook)
			return
		}

		// 关闭控制台输出,直接SetOutPut
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     false,
		})
		logrus.SetOutput(writer)
	})
}
