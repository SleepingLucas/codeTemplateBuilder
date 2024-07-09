package CreateTemplate

// CreateTemplate 创建代码模板接口
type CreateTemplate interface {
	CreateMain() (path string, err error) // 创建代码文件
	CreateTest() (path string, err error) // 创建测试文件
}
