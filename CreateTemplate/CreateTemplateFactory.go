package CreateTemplate

import (
	"github.com/SleepingLucas/ctb/CreateTemplate/impl"
)

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
