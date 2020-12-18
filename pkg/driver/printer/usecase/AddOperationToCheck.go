package printerUsecase

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"golang.org/x/text/encoding/charmap"
)

func renderStatus(status byte) string {
	const (
		sendData          = 1
		openedShift       = 2
		openedShift24hEnd = 3
		closedShift       = 4
		openedDocument    = 8
	)

	switch status {
	case sendData:
		return "отправка данных, 1"
	case openedShift:
		return "открытая смена, 24 часа не кончились , 2"
	case openedShift24hEnd:
		return "открытая смена, 24 часа не кончились , 3"
	case closedShift:
		return "закрытая смена"
	case openedDocument:
		return "открытый документ"
	default:
		return "неизвестный статус"
	}

}

func (p *printerUsecase) AddOperationToCheck(op models.Operation) error {
	p.logger.Debug("Send command AddOperationToCheck")

	// TODO: Проверка на на открытую смену перед началом чека, так как происходит
	// автоматическое отрытие смены без передачи инн кассира

	switch status := p.ReadShortStatus(); status {
	case models.OpenedShift,
		models.OpenedCheckIncome, models.OpenedCheckExpense, models.OpenedCheckReturnIncome,
		models.OpenedCheckReturnExpence, models.OpenedCheckNonFiscal:
		p.logger.Info("статус: ", renderStatus(status))
	default:
		p.logger.Info("Нельзя добавлять, статус: ", renderStatus(status))
		return errors.New(fmt.Sprint("Нельзя добавлять, статус: ", renderStatus(status)))
	}

	buf, cmdLen := p.createCommandBuffer(models.OperationV2, p.password)

	// Запись типа операции
	buf.WriteByte(op.Type)

	// Запись количества товара
	// Количество записывается в миллиграммах
	amount, err := p.intToBytesWithLen(op.Amount*consts.Milligram, 6)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	//p.logger.Debug("amount:\n", hex.Dump(amount))
	buf.Write(amount)

	// запись цены товара
	// цена записывается в копейках
	price, err := p.intToBytesWithLen(op.Price, 5) // одна копейка
	if err != nil {
		p.logger.Error(err)
		return err
	}

	buf.Write(price)

	// запись суммы товара
	// Сумма записывается в копейках
	summ, err := p.intToBytesWithLen(op.Sum, 5) // две копейки
	if err != nil {
		p.logger.Error(err)
		return err
	}
	buf.Write(summ)

	// Запись налогов на товар
	// Налог записывается в копейках
	//tax, err := intToBytesWithLen(0, 5)
	//if err != nil {
	//	p.logger.Fatal(err)
	//}
	buf.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff}) // если нет налога надо отправлять 0xff*6
	//buf.Write(tax)

	// Запись налоговой ставки
	buf.WriteByte(consts.VAT0)
	// Запись номера отдела
	buf.WriteByte(1)

	// Запись признака способа рассчета
	buf.WriteByte(consts.FullPayment)

	// Запись признака предмета рассчета
	buf.WriteByte(op.Subject)

	// Запись название товара 0 - 128 байт строка
	// кодировка win1251
	str, err := charmap.Windows1251.NewEncoder().String(op.Name)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	// создаем массив с длинной 128 байт
	rStrBytes := make([]byte, 128)
	copy(rStrBytes, str)
	buf.Write(rStrBytes[:128]) // записываем только первые 128 байт

	p.logger.Debug("длинна сообщения в байтах: ", buf.Len())
	p.logger.Debug("\n", hex.Dump(buf.Bytes()))
	p.logger.Debug("cmdlen: ", cmdLen)
	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Error(err)
		return err
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Error(err)
		return err
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

	return nil
}
