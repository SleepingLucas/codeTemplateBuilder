package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SleepingLucas/ctb/subcmd"
)

func main() {
	if len(os.Args) == 1 { // 如果没有参数
		flag.Usage()
		return
	}

	subCmd := subcmd.SubCmdFactory(os.Args[1])
	fmt.Println(os.Args[1])

	if err := subCmd.Exec(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
