package impl

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// CFTemplate codeforces 题目模板
type CFTemplate struct {
	ProblemName string // 题目名 例如 1840D
	URL         string // 题目链接
}

// isExist 判断文件是否存在
func (cf CFTemplate) isExist(suf string) (path string, exist bool) {
	// 解析字符串
	var contest int
	var problem string
	fmt.Sscanf(cf.ProblemName, "%d%s", &contest, &problem)
	path = fmt.Sprintf("%d_%s%s.go", contest, problem, suf)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return path, false
		}
		panic(err)
	}
	return path, true
}

// CreateMain 创建代码文件
func (cf CFTemplate) CreateMain() (codePath string, err error) {
	// 判断文件是否存在
	codePath, exist := cf.isExist("")
	if exist {
		fmt.Println("代码文件已存在，是否覆盖？(Y/n)")
		s := "Y"
		fmt.Scanln(&s)
		if strings.ToLower(s) != "y" {
			return "", nil
		} else {
			// 删除文件
			err := os.Remove(codePath)
			if err != nil {
				return "", err
			}
		}
	}

	file, err := os.OpenFile(codePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.WriteString(fmt.Sprintf(`package main

import (
	"bufio"
	. "fmt"
	"io"
	"os"
)

func cf%s(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	
}

func main() { cf%s(bufio.NewReader(os.Stdin), os.Stdout) }
`, cf.ProblemName, cf.ProblemName))

	return
}

// CreateTest 创建测试文件
func (cf CFTemplate) CreateTest() (testPath string, err error) {
	// 判断文件是否存在
	testPath, exist := cf.isExist("_test")
	if exist {
		fmt.Println("测试文件已存在，是否覆盖？(Y/n)")
		s := "Y"
		fmt.Scanln(&s)
		if strings.ToLower(s) != "y" {
			return "", nil
		} else {
			// 删除文件
			err := os.Remove(testPath)
			if err != nil {
				return "", err
			}
		}
	}

	file, err := os.OpenFile(testPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	testTemplatef := `package main

import (
	"testing"

	"github.com/EndlessCheng/codeforces-go/main/testutil"
)

func Test_cf%s(t *testing.T) {
	testCases := [][2]string{
%s
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf%s)
}`

	if cf.URL == "" {
		_, err = writer.WriteString(fmt.Sprintf(testTemplatef, cf.ProblemName,
			`		{
			`+"`在此键入样例输入`,\n"+
				"			`在此键入样例输出`,"+`
		},`, cf.ProblemName))
	} else {
		// 爬取题目样例
		inputs, outputs := cf.crawler(cf.URL)
		// off := 3 // 偏移量

		var totalBuilder strings.Builder

		// 组合测试用例
		for i := 0; i < len(inputs); i++ {
			var builder strings.Builder
			builder.WriteString("		{\n")

			// 写入输入用例
			builder.WriteString("			`")
			// 先读入第一行
			line, err := inputs[i].ReadString('\n')
			builder.WriteString(strings.TrimSpace(line))
			for {
				line, err = inputs[i].ReadString('\n')
				if err != nil {
					break
				}
				builder.WriteString("\n")
				// for range off {
				// 	builder.WriteString("\t")
				// }
				builder.WriteString(strings.TrimSpace(line))
			}
			builder.WriteString("`,\n")

			// 写入输出用例
			builder.WriteString("			`")
			// 先读入第一行
			line, err = outputs[i].ReadString('\n')
			builder.WriteString(strings.TrimSpace(line))
			for {
				line, err = outputs[i].ReadString('\n')
				if err != nil {
					break
				}
				builder.WriteString("\n")
				// for range off {
				// 	builder.WriteString("\t")
				// }
				builder.WriteString(strings.TrimSpace(line))
			}
			builder.WriteString("`,\n		},\n")

			totalBuilder.WriteString(builder.String())
		}

		_, err = writer.WriteString(
			fmt.Sprintf(testTemplatef,
				cf.ProblemName,
				strings.TrimSuffix(totalBuilder.String(), "\n"),
				cf.ProblemName,
			),
		)
	}

	return
}

// crawler 爬取 codeforces 题目样例
// 返回 inputs, outputs
func (cf CFTemplate) crawler(url string) (inputs, outputs []bytes.Buffer) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	// 创建 input、output map
	var inputmp, outputmp sync.Map
	n := 0

	var wg sync.WaitGroup

	doc.Find(".input").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var buf bytes.Buffer

			s.Find("pre div").Each(func(i int, s *goquery.Selection) {
				buf.WriteString(s.Text() + "\n")
			})
			if buf.Len() == 0 {
				s.Find("pre").Each(func(i int, s *goquery.Selection) {
					html, _ := s.Html()
					buf.WriteString(strings.ReplaceAll(html, "<br/>", "\n") + "\n")
				})
			}

			inputmp.Store(id, buf)
			n++

			// fmt.Printf("样例输入 %d 爬取成功\n", id+1)
		}(i)
	})

	doc.Find(".output").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var buf bytes.Buffer

			s.Find("pre").Each(func(i int, s *goquery.Selection) {
				buf.WriteString(s.Text())
			})

			outputmp.Store(id, buf)

			// fmt.Printf("样例输出 %d 爬取成功\n", id+1)
		}(i)
	})

	wg.Wait()
	// fmt.Println("完成爬取")

	for i := 0; i < n; i++ {
		if val, ok := inputmp.Load(i); ok {
			inputs = append(inputs, val.(bytes.Buffer))
		}
		if val, ok := outputmp.Load(i); ok {
			outputs = append(outputs, val.(bytes.Buffer))
		}
	}

	return
}
