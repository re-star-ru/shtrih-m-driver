package http

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/fess932/shtrih-m-driver/pkg/driver/printer"
	"github.com/kataras/iris/v12"
)

type PrinterHandler struct {
	password       uint32
	PrinterUsecase printer.Usecase
}

func NewPrinterHandler(party iris.Party, usecase printer.Usecase, host, password uint32) {
	handler := &PrinterHandler{
		password:       password,
		PrinterUsecase: usecase,
	}

	party.Post("/print_check", handler.PrintCheck)
}

func (p *PrinterHandler) PrintCheck(ctx iris.Context) {
	var chk models.CheckPackage

	if err := ctx.ReadJSON(&chk); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	for _, op := range chk.Operations {
		p.PrinterUsecase.AddOperationToCheck(op)
	}

	p.PrinterUsecase.CloseCheck(chk)
}
