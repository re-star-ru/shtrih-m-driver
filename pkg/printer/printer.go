package printer

import (
	"github.com/fess932/shtrih-m-driver/internal/logger"
	"github.com/fess932/shtrih-m-driver/internal/shtrih"
)

func NewPrinter(logger logger.Logger, host string, password uint32) *shtrih.Printer {
	return shtrih.NewPrinter(logger, host, password)
}
