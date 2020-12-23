package printerUsecase

import "github.com/fess932/shtrih-m-driver/pkg/driver/models"

func (p *printerUsecase) Print(chk models.CheckPackage) error {
	switch s := p.ReadShortStatus(); s {
	case models.OpenedCheckIncome, models.OpenedCheckExpense:
		if err := p.CancellationOpenedCheck(); err != nil {
			p.logger.Debug(err)
			return err
		}
	}

	// добавляем операции, создается чек
	for _, v := range chk.Operations {
		if err := p.AddOperationToCheck(v); err != nil {
			// TODO: Очистка операций в кеше перед выходом
			p.logger.Error(err)
			return err
		}
	}

	// закрываем чек
	return p.CloseCheck(chk)
}
