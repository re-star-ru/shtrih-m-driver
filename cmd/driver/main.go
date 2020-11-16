package main

import (
	"log"

	"github.com/fess932/shtrih-m-driver/pkg/driver/client/usecase/tcp"

	printerUsecase "github.com/fess932/shtrih-m-driver/pkg/driver/printer/usecase"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	//host := "10.51.0.71:7778"
	//password := uint32(30)

	host := "fake"
	password := uint32(0000)

	//c := emulator.NewClientUsecase(host, logger)
	//p := printerUsecase.NewPrinterUsecase(logger, c, password)
	//p.ReadShortStatus()

	host = "10.51.0.71:7778"
	password = uint32(30)

	c := tcp.NewClientUsecase(host, logger)
	p := printerUsecase.NewPrinterUsecase(logger, c, password)
	p.ReadShortStatus()

	//p.AddOperationToCheck(models.Operation{
	//	Type:    0,
	//	Amount:  0,
	//	Price:   0,
	//	Sum:     0,
	//	Subject: 0,
	//	Name:    "",
	//})
	//
	//p.CloseCheck(models.CheckPackage{
	//	Cash:       0,
	//	Casheless:  0,
	//	TaxSystem:  0,
	//	BottomLine: "",
	//})

	//app := iris.New()
	//party := app.Party("/printer")

	//http.NewPrinterHandler(party, p, host, password)
	//
	//if err := app.Listen(":8080"); err != nil {
	//	logger.Error(err)
	//}
}
