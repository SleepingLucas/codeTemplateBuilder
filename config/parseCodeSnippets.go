package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"regexp"
)

// PrintToConsole vscode 用户代码片段的 Print to console 部分
type PrintToConsole struct {
	Scope       string   `json:"scope"`
	Prefix      string   `json:"prefix"`
	Body        []string `json:"body"`
	Description string   `json:"description"`
}

// CodeSnippet vscode 用户代码片段
type CodeSnippet struct {
	PrintToConsole `json:"Print to console"`
}

var (
	tabRe = regexp.MustCompile(`"\t*.*"`)          // 防止制表符导致的解析错误
	tplRe = regexp.MustCompile(`\$([A-Za-z0-9]+)`) // 替换模板中的变量

	stringRe  = regexp.MustCompile(`"([^"\\]|\\.)*"`) // 匹配字符串
	commentRe = regexp.MustCompile(`//.*\n`)
)

// RemoveComments 移除代码片段中的 json 注释
func removeComments(input []byte) []byte {
	// 将字符串中的注释替换为不匹配的字符串
	strippedString := stringRe.ReplaceAllFunc(input, func(s []byte) []byte {
		return bytes.ReplaceAll(s, []byte("//"), []byte("/*/"))
	})

	// 删除匹配的注释
	result := commentRe.ReplaceAll(strippedString, []byte(""))

	// 恢复原样
	result = bytes.ReplaceAll(result, []byte("/*/"), []byte("//"))

	// 返回处理后的字符串
	return result
}

// ParseVsCodeSnippet 解析 vscode 代码片段
func ParseVsCodeSnippet(path string) (code []string, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// 读取并解析文件内容
	byteCodeValue, _ := io.ReadAll(jsonFile)
	newCodeBytes := tabRe.ReplaceAllFunc(byteCodeValue, func(s []byte) []byte {
		return bytes.ReplaceAll(s, []byte("\t"), []byte("\\t"))
	})
	newCodeBytes = removeComments(newCodeBytes)
	newCodeBytes = bytes.ReplaceAll(newCodeBytes, []byte("$0"), []byte(""))
	newCodeBytes = tplRe.ReplaceAll(newCodeBytes, []byte("{{ .d$1 }}"))

	codeSnippet := new(CodeSnippet)
	err = json.Unmarshal(newCodeBytes, codeSnippet)
	if err != nil {
		return nil, err
	}

	return codeSnippet.Body, nil
}
