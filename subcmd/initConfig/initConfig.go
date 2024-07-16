package initConfig

import (
	"flag"

	"github.com/SleepingLucas/ctb/config"
)

type InitConfig struct {
	initFlagSet *flag.FlagSet // flag 集合
	reset       bool          // 重置配置文件
	cfCodePath  string        // codeforces 代码片段路径
	cfTestPath  string        // codeforces 测试片段路径
}

func (i *InitConfig) Init() error {
	i.initFlagSet = flag.NewFlagSet("init", flag.ExitOnError)

	i.initFlagSet.BoolVar(&i.reset, "reset", false, "重置配置文件")
	i.initFlagSet.StringVar(&i.cfCodePath, "cfCode", "", "codeforces 代码片段路径")
	i.initFlagSet.StringVar(&i.cfTestPath, "cfTest", "", "codeforces 测试片段路径")

	return nil
}

func (i *InitConfig) Run(args []string) (err error) {
	if err = i.initFlagSet.Parse(args); err != nil {
		return err
	}

	if i.reset {
		// 重置配置文件为默认配置
		config.WriteDefaultConfig(config.GetConfigPath())
		return nil
	}

	// 初始化配置文件
	if err = config.InitConfig(); err != nil {
		return err
	}
	var isChange bool

	if i.cfCodePath != "" {
		codeTemplate, err := config.ParseVsCodeSnippet(i.cfCodePath)
		if err != nil {
			return err
		}
		config.Conf.Templates.Codeforces.Code = codeTemplate
		isChange = true
	}

	if i.cfTestPath != "" {
		testTemplate, err := config.ParseVsCodeSnippet(i.cfTestPath)
		if err != nil {
			return err
		}
		config.Conf.Templates.Codeforces.Test = testTemplate
		isChange = true
	}

	if isChange {
		// 写入配置文件
		if err = config.OverrideConfig(config.GetConfigPath(), *config.Conf); err != nil {
			return
		}
	}

	return nil
}

func (i *InitConfig) Exec(args []string) error {
	if err := i.Init(); err != nil {
		return err
	}

	return i.Run(args[2:])
}
