package createTemplate

import (
	"github.com/SleepingLucas/ctb/CreateTemplate/impl"
)

func CreateTemplateFactory(templateType, problem, url string) CreateTemplate {
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
