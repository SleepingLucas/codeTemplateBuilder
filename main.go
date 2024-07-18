package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SleepingLucas/ctb/subcmd/ctb"

	"github.com/SleepingLucas/ctb/subcmd/initConfig"

	"github.com/SleepingLucas/ctb/subcmd"
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: ctb <command> [arguments]")
		fmt.Println("The commands are:")
		fmt.Println(" 空	生成代码片段")
		ctb := new(ctb.Ctb)
		ctb.Init()
		flag.PrintDefaults()
		fmt.Println()

		fmt.Println(" init	初始化配置文件")
		initConfig := &initConfig.InitConfig{}
		initConfig.Init()
		initConfig.InitFlagSet.PrintDefaults()
	}
}

func main() {
	if len(os.Args) == 1 { // 如果没有参数
		flag.Usage()
		return
	}

	subCmd := subcmd.Factory(os.Args[1])

	if err := subCmd.Exec(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
