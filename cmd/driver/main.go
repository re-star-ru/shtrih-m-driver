package main

import (
	"log"

	"github.com/fess932/shtrih-m-driver/pkg/driver"
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

	host := "10.51.0.71:7778"
	password := uint32(30)

	p := driver.NewPrinter(logger, host, password)
	//driver.SellOperationV2()

	if err := p.TLVWriteCashierINN("263209745357"); err != nil {
		logger.Fatal(err)
	}

	//driver.CloseCheckV2()

	//p.ReadShortStatus()
	//
	//driver.FnReadStatus()

	//driver.OpenCheck()
	//driver.AddSale(2000, 100)
	//driver.CloseCheck()
	//
	//driver.ReadShortStatus()

	//driver.PrintReportWithoutClearing()
	//driver.PrintCheck()

	//driver.WriteTable(tables.TableCashier, 14, 2, "Оператор14")
	//driver.ReadFieldInfo(shtrih.SmfpTableCashier, 1)

	//testTLV(logger)
}

//func testTLV(log *zap.SugaredLogger) {
//	data := driver.New(consts.FNNumber, 15)
//	log.Debug(data)
//
//	log.Debug(binary.LittleEndian.Uint16(data[:2]))
//	log.Debug(binary.LittleEndian.Uint16(data[2:4]))
//
//}
