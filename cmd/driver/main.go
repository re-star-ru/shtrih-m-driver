package main

import (
	"log"
	"shtrih-drv/internal/shtrih"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func createLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, "current.log")
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

	host := "10.51.0.71:7778"
	password := uint32(30)

	printer := shtrih.NewPrinter(logger, host, password)
	//printer.FnReadStatus()

	//printer.ReadShortStatus()
	//printer.PrintReportWithoutClearing()
	//printer.PrintReportWithoutClearing()

	printer.PrintCheck()
}
