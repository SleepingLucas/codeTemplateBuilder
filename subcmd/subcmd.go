package subcmd

import (
	"github.com/SleepingLucas/ctb/subcmd/ctb"
	"github.com/SleepingLucas/ctb/subcmd/initConfig"
)

type SubCmd interface {
	Init() error              // 对子命令进行初始化参数的绑定
	Run(args []string) error  // 解析参数并执行
	Exec(args []string) error // 执行子命令
	PrintDefaults()           // 打印帮助信息
}

var (
	initC = new(initConfig.InitConfig)
	ctbC  = new(ctb.Ctb)
)

// Factory 子命令工厂
func Factory(cmd string) SubCmd {
	switch cmd {
	case "init":
		return initC
	default:
		return ctbC
	}
}
