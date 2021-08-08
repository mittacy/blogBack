package transform

import (
	"github.com/mittacy/blogBack/pkg/logger"
)

type Email struct {
	logger *logger.CustomLogger
}

func NewEmail(customLogger *logger.CustomLogger) Email {
	return Email{logger: customLogger}
}

