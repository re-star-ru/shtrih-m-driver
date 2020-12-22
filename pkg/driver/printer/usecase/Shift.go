package printerUsecase

import (
	"errors"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

// shift open

func (p *printerUsecase) OpenShift(c models.Cashier) error {

	// проверка статуса
	// если статус смена открыта то вернуть ошибку
	switch status := p.ReadShortStatus(); status {
	case models.OpenedShift:
		return errors.New("cмена уже открыта")
	}

	// записать имя кассира,
	//p.writeCashierName(c.Name)

	// начать открытие смен
	if err := p.startingShiftOpening(); err != nil {
		return err
	}

	// записать инн кассира
	if err := p.writeCashierINN(c.INN); err != nil {
		return err
	}

	// открыть смену
	if err := p.shiftOpening(); err != nil {
		return err
	}

	return nil
}

func (p *printerUsecase) startingShiftOpening() error {
	p.logger.Debug("Send command startingShiftOpening")

	buf, cmdLen := p.createCommandBuffer(models.StartOpenShift, p.password)

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

func (p *printerUsecase) shiftOpening() error {
	p.logger.Debug("Send command shiftOpening")

	buf, cmdLen := p.createCommandBuffer(models.OpenShift, p.password)

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

// shift close

func (p *printerUsecase) CloseShift() error {
	// проверка статуса
	// если статус смена не открыта то вернуть ошибку
	if status := p.ReadShortStatus(); status != models.OpenedShift {
		return errors.New("смена не открыта")
	}

	if err := p.closingShift(); err != nil {
		return err
	}

	return nil
}

func (p *printerUsecase) closingShift() error {
	p.logger.Debug("Send command ZReport")

	buf, cmdLen := p.createCommandBuffer(models.ZReport, p.password)

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
