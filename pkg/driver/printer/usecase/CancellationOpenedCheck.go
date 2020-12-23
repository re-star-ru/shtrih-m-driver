package printerUsecase

import (
	"errors"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

func (p *printerUsecase) CancellationOpenedCheck() error {
	p.logger.Debug("Send command CancellationOpenedCheck")

	switch status := p.ReadShortStatus(); status {
	case models.OpenedCheckIncome, models.OpenedCheckExpense, models.OpenedCheckReturnIncome,
		models.OpenedCheckReturnExpence, models.OpenedCheckNonFiscal:
		p.logger.Info("статус: ", status)
	default:
		err := errors.New("нет чека для отмены")
		p.logger.Debug(err)
		return err
	}

	buf, cmdLen := p.createCommandBuffer(models.CancellationCheck, p.password)

	rFrame, err := p.send(buf.Bytes(), cmdLen)
	if err != nil {
		p.logger.Debug(err)
		return err
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Debug(err)
		return err
	}
	return nil
}
