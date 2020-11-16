package emulator

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/client"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/fess932/shtrih-m-driver/pkg/logger"
)

type Usecase struct {
	host   string
	logger logger.Logger
}

func NewClientUsecase(host string, logger logger.Logger) client.Usecase {
	return &Usecase{host: host, logger: logger}
}

func (u *Usecase) Send(frame []byte, cmdLen int) (*models.Frame, error) {
	u.parse(frame)

	return &models.Frame{}, nil
}

// parse Парсит frame и определяет к какой операции он относится и какие данные необходимо вернуть
func (u *Usecase) parse(frame []byte) {
	u.logger.Debug("debug!")
	u.logger.Debug(frame)
}
