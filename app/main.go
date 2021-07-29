package main

import (
	"log"
	"net/http"

	"github.com/fess932/shtrih-m-driver/app/kkt"

	"github.com/go-chi/chi/v5"

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

	r := chi.NewRouter()
	confRouter(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func confRouter(r *chi.Mux) {
	k := kkt.New("13")

	r.Get("/open", func(w http.ResponseWriter, r *http.Request) {
		if err := k.OpenShift(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	r.Get("/close", func(w http.ResponseWriter, r *http.Request) {
		if err := k.CloseShift(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

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
