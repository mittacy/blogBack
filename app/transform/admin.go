package transform

import (
	"github.com/mittacy/blogBack/pkg/logger"
)

type Admin struct {
	logger *logger.CustomLogger
}

func NewAdmin(customLogger *logger.CustomLogger) Admin {
	return Admin{logger: customLogger}
}

