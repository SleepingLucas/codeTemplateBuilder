package CreateTemplate

import "github.com/SleepingLucas/ctb/subcmd/ctb/CreateTemplate/impl"

// CreateTemplate 创建代码模板接口
type CreateTemplate interface {
	CreateMain() (path string, err error) // 创建代码文件
	CreateTest() (path string, err error) // 创建测试文件
}

// Factory 创建代码模板工厂
func Factory(templateType, problem, url string) CreateTemplate {
	switch templateType {
	case "cf":
		if url != "" {
			return impl.CFTemplate{ProblemName: problem, URL: url}
		} else {
			return impl.CFTemplate{ProblemName: problem}
		}
	default:
		panic("unknown template type")
	}
}
