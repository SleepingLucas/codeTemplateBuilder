package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
)

var Conf = new(Config)

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
	// 配置文件的路径
	configFilePath := GetConfigPath()

	// 检查配置文件是否存在
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// 不存在则创建并写入默认配置
		WriteDefaultConfig(configFilePath)
	} else if err != nil {
		panic(err)
	}

	// 读取配置文件
	if err = UnmarshalConfig(configFilePath); err != nil {
		return err
	}

	// 判断是否有空配置，有则写入默认配置
	val := reflect.ValueOf(Conf).Elem()
	defaultVal := reflect.ValueOf(defaultConfig)
	flag := false
	// traverseAndCheckEmpty 遍历配置文件并检查是否有空配置
	var traverseAndCheckEmpty func(nowField reflect.Value, defaultField reflect.Value)
	traverseAndCheckEmpty = func(nowField reflect.Value, defaultField reflect.Value) {
		for i := 0; i < nowField.NumField(); i++ {
			field := nowField.Field(i)
			// 检查字段是否为结构体，如果是，则递归遍历
			if field.Kind() == reflect.Struct {
				traverseAndCheckEmpty(field, defaultField.Field(i))
			} else {
				if field.IsZero() {
					// 将默认配置 defaultConfig 的对应字段写入 Conf
					field.Set(defaultField.Field(i))
					flag = true
				}
			}
		}
	}
	traverseAndCheckEmpty(val, defaultVal)
	if flag {
		// 写入配置
		if err := OverrideConfig(configFilePath, *Conf); err != nil {
			panic(err)
		}
	}

	return nil
}

// OverrideConfig 覆盖写入配置文件
func OverrideConfig(configFilePath string, config Config) error {
	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	// 获取用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// 配置文件的路径
	configFilePath := filepath.Join(homeDir, ".ctbconfig")

	return configFilePath
}

// UnmarshalConfig 读取配置文件到 Conf
func UnmarshalConfig(configFilePath string) (err error) {
	// 读取配置文件
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("json")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
		return
	}

	// 反序列化配置文件
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("unmarshal to Conf failed, err:%v\n", err)
		return
	}

	return
}
