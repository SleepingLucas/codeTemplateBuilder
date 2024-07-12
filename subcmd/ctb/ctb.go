package ctb

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"sync"

	"github.com/SleepingLucas/ctb/config"

	"github.com/SleepingLucas/ctb/subcmd/ctb/CreateTemplate"

	"github.com/pkg/errors"
)

var (
	cfURL1 = regexp.MustCompile(`^https://codeforces.com/contest/(\d+)/problem/([a-zA-Z])$`)
	cfURL2 = regexp.MustCompile(`^https://codeforces.com/problemset/problem/(\d+)/([a-zA-Z])$`)

	ErrorNoProblemNameOrURL = errors.New("请输入题目名或题目网址\n\t如: -p=1840D\n\t如: -url=https://codeforces.com/contest/1840/problem/D")
	ErrorProblemNameFormat  = errors.New("题目名格式错误\n\t如: -p=1840D")
	ErrorURL                = errors.New("题目链接错误")
)

type Ctb struct {
	templateType string // 模板类型
	problemName  string // 题目名
	testOnly     bool   // 只生成测试文件
	codeOnly     bool   // 只生成代码文件
	url          string // 题目链接
}

func (c *Ctb) Init() error {
	flag.StringVar(&c.templateType, "type", "cf", "模板类型")
	flag.StringVar(&c.problemName, "problem", "", "题目名") // 例如 1840D
	flag.StringVar(&c.problemName, "p", "", "题目名(shortcut)")
	flag.BoolVar(&c.testOnly, "test", false, "只生成测试文件")
	flag.BoolVar(&c.codeOnly, "code", false, "只生成代码文件")
	flag.StringVar(&c.url, "url", "", "题目链接")

	return nil
}

func (c *Ctb) Run(args []string) error {
	flag.Parse()

	fmt.Printf("ctb: %+v\n", c)

	// 参数校验
	if c.url == "" {
		if c.problemName == "" {
			return ErrorNoProblemNameOrURL
		}
		if ok, err := regexp.MatchString(`^\d+[A-Z]$`, c.problemName); !ok || err != nil {
			return ErrorProblemNameFormat
		}
	} else {
		c.problemName = GetProblemName(c.url)
		if c.problemName == "" {
			return ErrorURL
		}
	}

	factory := CreateTemplate.Factory(c.templateType, c.problemName, c.url)

	var wg sync.WaitGroup

	if !c.testOnly {
		// 生成代码文件
		wg.Add(1)
		go func() {
			defer wg.Done()
			codePath, err := factory.CreateMain()
			if err != nil {
				panic(err)
			}

			// 执行命令以在vscode中打开文件：code codePath
			cmd := exec.Command("code", codePath)
			_ = cmd.Run()
		}()
	}

	if !c.codeOnly {
		// 生成测试文件
		wg.Add(1)
		go func() {
			defer wg.Done()
			testPath, err := factory.CreateTest()
			if err != nil {
				panic(err)
			}

			cmd := exec.Command("code", testPath)
			_ = cmd.Run()
		}()
	}

	wg.Wait()

	return nil
}

func (c *Ctb) Exec(args []string) error {
	configFilePath := config.GetConfigPath() // 获取配置文件路径

	// 读取配置文件
	if err := config.UnmarshalConfig(configFilePath); err != nil {
		return err
	}

	if err := c.Init(); err != nil {
		return err
	}

	if err := c.Run(args); err != nil {
		return err
	}

	return nil
}

// GetProblemName 从题目链接中获取题目名
func GetProblemName(url string) string {
	// 目前有两种链接格式
	// https://codeforces.com/contest/1926/problem/G
	// https://codeforces.com/problemset/problem/1759/E

	// 解析第一种链接
	match := cfURL1.FindStringSubmatch(url)
	if len(match) == 3 {
		return fmt.Sprintf("%s%s", match[1], match[2])
	}

	// 解析第二种链接
	match = cfURL2.FindStringSubmatch(url)
	if len(match) == 3 {
		return fmt.Sprintf("%s%s", match[1], match[2])
	}

	return ""
}
