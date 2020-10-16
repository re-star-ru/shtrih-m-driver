package main

import (
	"encoding/binary"
	"log"

	"github.com/fess932/shtrih-m-driver/internal/shtrih/TLV"
	"github.com/fess932/shtrih-m-driver/pkg/printer"

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

	p := printer.NewPrinter(logger, host, password)
	//printer.SellOperationV2()
	//printer.CloseCheckV2()

	p.ReadShortStatus()
	//
	//printer.FnReadStatus()

	//printer.OpenCheck()
	//printer.AddSale(2000, 100)
	//printer.CloseCheck()
	//
	//printer.ReadShortStatus()

	//printer.PrintReportWithoutClearing()
	//printer.PrintCheck()

	//printer.WriteTable(tables.TableCashier, 14, 2, "Оператор14")
	//printer.ReadFieldInfo(shtrih.SmfpTableCashier, 1)

	//testTLV(logger)
}
func testTLV(log *zap.SugaredLogger) {
	data := TLV.New(TLV.FNNumber, 15)
	log.Debug(data)

	log.Debug(binary.LittleEndian.Uint16(data[:2]))
	log.Debug(binary.LittleEndian.Uint16(data[2:4]))

}
