package config

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/HideInBush7/go_im/pkg/util"
	"github.com/spf13/viper"
)

func init() {
	Init()
}

var once sync.Once

func Init() {
	once.Do(func() {
		var confFile string

		// 环境变量env=dev开发状态下 配置文件为相对路径下的config.yaml #方便测试
		if env, ok := os.LookupEnv(`env`); ok && env == `dev` {
			confFile = "./config.yaml"
		} else {
			// 配置文件可以通过参数--config输入,默认为执行文件的同个目录,直接go run的话会在/tmp下的一个临时目录中
			flag.StringVar(&confFile, `config`, filepath.Join(util.GetExecDir(), "config.yaml"), `set config file`)
			flag.Parse()
		}
		viper.SetConfigFile(confFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	})
}

func Register(key string, in interface{}) {
	err := viper.UnmarshalKey(key, in)
	if err != nil {
		panic(err)
	}
}
