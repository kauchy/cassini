// Package commands 实现命令行应用基本命令
//
// 定义 Root 命令（默认命令）以实现默认功能（显示帮助信息），并实现预处理功能。
//
// 定义 $> cassini start 命令(服务启动命令)以实现服务启动，并根据配置运行服务。
//
// 定义 $> cassini mock 命令(Mock服务启动命令)以实现Mock服务启动，并根据配置运行Mock服务，以便于进行服务相关测试。
//
// 定义 $> cassini wsclient 启动WebSocket客户端，以监听服务端交易事件，进行相关测试。
package commands

import (
	"context"
	"os"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/spf13/cobra"
)

var conf string
var logConf string

// Runner 通过配置数据执行方法，返回运行过程中出现的错误，如果返回空则代表运行成功。
type Runner func(conf *config.Config) (context.CancelFunc, error)

// NewRootCommand 创建 root/默认 命令
//
// 实现默认功能，显示帮助信息，预处理配置初始化，日志配置初始化。
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "cassini",
		Short: "relay between blockchains",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			_, err = config.LoadConfig(conf)
			if err != nil {
				log.Error("Run root command error: ", err.Error())
				return
			}
			log.Debug("Init config: ", conf)
			return
		},
	}

	return root
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run(config.GetConfig())
	if err != nil {
		log.Error("Run command error: ", err.Error())
		return err
	}
	if isKeepRunning {
		common.KeepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}
