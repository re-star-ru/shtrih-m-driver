package main

import (
	"log"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

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

	// TODO:  передача инн кассира перед первой операцией
	// 1021, кассир
	// printer.writeCashierName(params.Parameters.CashierName);

	p.AddOperationToCheck(models.Operation{
		Type:    consts.Income,
		Amount:  1,
		Price:   1,
		Sum:     1,
		Subject: consts.Service,
		Name:    "Ремонт стартера тест",
	})

	p.AddOperationToCheck(models.Operation{
		Type:    consts.Income,
		Amount:  1,
		Price:   1,
		Sum:     1,
		Subject: consts.Service,
		Name:    "Ремонт стартера тест",
	})
	//

	////возврат прихода
	//p.AddOperationToCheck(models.Operation{
	//	Type:    consts.ReturnIncome,
	//	Amount:  1,
	//	Price:   1,
	//	Sum:     1,
	//	Subject: consts.Service,
	//	Name:    "Ремонт стартера тест",
	//})
	//
	//// возврат прихода
	//p.AddOperationToCheck(models.Operation{
	//	Type:    consts.ReturnIncome,
	//	Amount:  1,
	//	Price:   1,
	//	Sum:     1,
	//	Subject: consts.Service,
	//	Name:    "Ремонт стартера тест",
	//})

	if err := p.WriteCashierINN("263209745357"); err != nil {
		logger.Error(err)
	}

	p.CloseCheck(models.CheckPackage{
		Operations: nil,
		Cash:       2,
		Casheless:  0,
		TaxSystem:  consts.ENVD,
		BottomLine: "Нижняя линия чека",
	}, true)
}
