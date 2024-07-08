package createTemplate

// CreateTemplate 创建代码模板接口
type CreateTemplate interface {
	CreateMain() error // 创建代码文件
	CreateTest() error // 创建测试文件
}
