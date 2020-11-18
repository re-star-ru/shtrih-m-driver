package printerUsecase

import "github.com/fess932/shtrih-m-driver/pkg/driver/models"

func (p *printerUsecase) CancellationOpenedCheck() {
	p.logger.Debug("Send command CancellationOpenedCheck")

	switch status := p.ReadShortStatus(); status {

	case models.OpenedCheckIncome, models.OpenedCheckExpense, models.OpenedCheckReturnIncome,
		models.OpenedCheckReturnExpence, models.OpenedCheckNonFiscal:
		p.logger.Info("статус: ", status)
	default:
		p.logger.Debug("Нет чека для отмены")
		return
	}

	buf, cmdLen := p.createCommandBuffer(models.CancellationCheck, p.password)

	rFrame, err := p.send(buf.Bytes(), cmdLen)
	if err != nil {
		p.logger.Error(err)
		return
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Error(err)
		return
	}

}
