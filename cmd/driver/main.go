package main

import (
	"log"
	"shtrih-drv/internal/fiscalprinter"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func main() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal()
	}
	slogger := logger.Sugar()

	slogger.Info("Shtih driver starting")
	//slogger.Debug("This is a DEBUG message")
	//slogger.Info("This is an INFO message")
	////slogger.Info("This is an INFO message with fields", "region", "us-west", "id", 2)
	//slogger.Warn("This is a WARN message")
	//slogger.Error("This is an ERROR message")

	printer := fiscalprinter.NewPrinterProtocol(slogger)
	err = printer.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
