package main

import (
	"fmt"
	"log"

	"github.com/fess932/shtrih-m-driver/app/kkt"

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

	k := kkt.New("13")

	if err := k.FSM.Event("open"); err != nil {
		log.Println(err)
	}

	if err := k.FSM.Event("close"); err != nil {
		fmt.Println(err)
	}

	//
	//host := "fake"
	//password := uint32(0000)
	//
	//
	//
	//c := emulator.NewClientUsecase(host, logger)
	////p := printerUsecase.NewPrinterUsecase(logger, c, password)
	////p.ReadShortStatus()
	//
	////host = "10.51.0.73:7778"
	////password = uint32(30)
	//
	////c := tcp.NewClientUsecase(host, time.Millisecond*5000, logger)
	//p := printerUsecase.NewPrinterUsecase(logger, c, password)
	//
	//if err := p.FNOpenedDocumentCancel(); err != nil {
	//	logger.Error(err)
	//}
}