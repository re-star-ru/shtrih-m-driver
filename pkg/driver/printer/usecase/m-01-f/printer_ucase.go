package m_01_f

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/fess932/shtrih-m-driver/pkg/driver/printer"
)

type printerUsecase struct {
}

func NewPrinterUsecase() printer.Usecase {
	return printerUsecase{}
}

func (p printerUsecase) OpenShift() {
	panic("implement me")
}

func (p printerUsecase) CloseShift() {
	panic("implement me")
}

func (p printerUsecase) AddOperationToCheck(op models.Operation) {
	panic("implement me")
}

func (p printerUsecase) CloseCheck() {
	panic("implement me")
}

func (p printerUsecase) SellOperationV2(op models.Operation) {
	panic("implement me")
}
