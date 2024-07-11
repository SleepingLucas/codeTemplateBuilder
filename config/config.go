package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Conf = new(Config)

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

// Template 代码片段模板
type Template struct {
	Code []string `json:"code"` // 代码模板
	Test []string `json:"test"` // 测试模板
}

// Templates 代码片段模板集合
type Templates struct {
	Codeforces Template `json:"codeforces"` // Codeforces 模板
}

// Config ctb 配置文件，存储代码片段，包括但不限于 CF 模板
type Config struct {
	Templates `json:"templates"`
}

func InitConfig() error {
	// 获取用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// 配置文件的路径
	configFilePath := filepath.Join(homeDir, ".ctbconfig")

	// 检查配置文件是否存在
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// 不存在则创建并写入默认配置
		writeDefaultConfig(configFilePath)
	} else if err != nil {
		panic(err)
	}

	// 读取配置文件
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 反序列化配置文件
	if err = viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v\n", err))
	}

	return nil
}

// WriteConfig 写入配置文件
func WriteConfig(file *os.File, config Config) error {
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}
