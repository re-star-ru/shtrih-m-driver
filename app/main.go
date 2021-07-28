package main

import (
	"log"

	"github.com/fess932/shtrih-m-driver/pkg/driver/client/usecase/emulator"

	printerUsecase "github.com/fess932/shtrih-m-driver/pkg/driver/printer/usecase"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0000")
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal()
	}

	slogger := logger.Sugar()
	slogger.Debug("log level: ", loggerConfig.Level.String())
	return slogger
}

func main() {
	logger := createLogger()
	logger.Info("Shtrih driver starting")

	host := "fake"
	password := uint32(0000)

	c := emulator.NewClientUsecase(host, logger)
	//p := printerUsecase.NewPrinterUsecase(logger, c, password)
	//p.ReadShortStatus()

	//host = "10.51.0.73:7778"
	//password = uint32(30)

	//c := tcp.NewClientUsecase(host, time.Millisecond*5000, logger)
	p := printerUsecase.NewPrinterUsecase(logger, c, password)

	if err := p.FNOpenedDocumentCancel(); err != nil {
		logger.Error(err)
	}
}
