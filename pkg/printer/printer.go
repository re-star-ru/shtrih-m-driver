package printer

import (
	"github.com/fess932/shtrih-m-driver/internal/shtrih"
	"github.com/fess932/shtrih-m-driver/pkg/logger"
)

func NewPrinter(logger logger.Logger, host string, password uint32) *shtrih.Printer {
	return shtrih.NewPrinter(logger, host, password)
}
