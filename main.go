package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/SleepingLucas/ctb/CreateTemplate"
)

var (
	templateType string // 模板类型
	problemName  string // 题目名
	testOnly     bool   // 只生成测试文件
	codeOnly     bool   // 只生成代码文件
	url          string // 题目链接

	cfURL1 = regexp.MustCompile(`^https://codeforces.com/contest/(\d+)/problem/([a-zA-Z])$`)
	cfURL2 = regexp.MustCompile(`^https://codeforces.com/problemset/problem/(\d+)/([a-zA-Z])$`)
)

func main() {
	flag.StringVar(&templateType, "type", "cf", "模板类型")
	flag.StringVar(&problemName, "problem", "", "题目名") // 例如 1840D
	flag.StringVar(&problemName, "p", "", "题目名(shortcut)")
	flag.BoolVar(&testOnly, "test", false, "只生成测试文件")
	flag.BoolVar(&codeOnly, "code", false, "只生成代码文件")
	flag.StringVar(&url, "url", "", "题目链接")

	flag.Parse()

	// 参数校验
	if url == "" {
		if problemName == "" {
			fmt.Println("请输入题目名\n\t如: -p=1840D")
			return
		}
		if ok, err := regexp.MatchString(`^\d+[A-Z]$`, problemName); !ok || err != nil {
			fmt.Println("题目名格式错误\n\t如: -p=1840D")
			return
		}
	} else {
		problemName = GetProblemName(url)
		if problemName == "" {
			fmt.Println("题目链接错误")
			return
		}
	}

	factory := CreateTemplate.CreateTemplateFactory(templateType, problemName, url)

	if !testOnly {
		// 生成代码文件
		codePath, err := factory.CreateMain()
		if err != nil {
			panic(err)
		}

		// 执行命令以在vscode中打开文件：code codePath
		cmd := exec.Command("code", codePath)
		_ = cmd.Run()
	}

	if !codeOnly {
		// 生成测试文件
		testPath, err := factory.CreateTest()
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("code", testPath)
		_ = cmd.Run()
	}
}

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
