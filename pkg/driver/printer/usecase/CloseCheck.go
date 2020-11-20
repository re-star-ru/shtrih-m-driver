package printerUsecase

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"golang.org/x/text/encoding/charmap"
)

func (p *printerUsecase) CloseCheck(chk models.CheckPackage, dontPrint bool) {
	p.logger.Debug("Send command CloseCheck")

	switch status := p.ReadShortStatus(); status {
	case models.OpenedCheckIncome, models.OpenedCheckExpense, models.OpenedCheckReturnIncome,
		models.OpenedCheckReturnExpence, models.OpenedCheckNonFiscal:
		p.logger.Info("статус: ", status)
	default:
		p.logger.Debug("Нет чека для закрытия")
		return
	}

	if dontPrint {
		p.DontPrintOneCheck() // не печатать чек если передан флаг dont print
	}

	if err := p.WriteCashierINN(chk.CashierINN); err != nil {
		// запись inn кассира после записи операций но до закрытия чека
		p.logger.Error(err)
		return
	}

	buf, cmdLen := p.createCommandBuffer(models.CloseCheckV2, p.password)
	p.logger.Debug("cmdlen:", cmdLen)

	// запись суммы наличных - типа оплаты 1
	cash, err := p.intToBytesWithLen(chk.Cash, 5)
	p.logger.Debug(len(cash), " длинна кеша")
	if err != nil {
		p.logger.Error(err)
		return
	}
	buf.Write(cash)

	// запись суммы типа оплаты 2 - безнал
	casheless, err := p.intToBytesWithLen(chk.Casheless, 5)
	if err != nil {
		p.logger.Error(err)
		return
	}
	buf.Write(casheless)

	// запись остальных с 3 по 16 видоов оплаты, длинна вида 5 байт, итого 70 байт
	buf.Write(make([]byte, 70))
	buf.WriteByte(0) // округление до рубля

	// запись налогов 6 * 5, итого 30 байт
	buf.Write(make([]byte, 30))

	buf.WriteByte(chk.TaxSystem) //система налогоообложения, биты а не байт

	// Запись нижней линии чека 64 байт строка
	str, err := charmap.Windows1251.NewEncoder().String(chk.BottomLine)
	if err != nil {
		p.logger.Error(err)
		return
	}
	rStrBytes := make([]byte, 64)
	copy(rStrBytes, str)

	buf.Write(rStrBytes)

	p.logger.Debug("len: ", buf.Len())

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Error(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))
}

type TLV struct {
	Tag   []byte
	Len   []byte
	Value []byte
}

func newTLV(Tag, Len uint16, value []byte) (TLV, error) {
	tlv := TLV{
		Tag:   make([]byte, 2),
		Len:   make([]byte, 2),
		Value: make([]byte, Len),
	}
	binary.LittleEndian.PutUint16(tlv.Tag, Tag) // код тега
	binary.LittleEndian.PutUint16(tlv.Len, Len) // длинна тега

	copy(tlv.Value, value) // значение тега

	if len(tlv.Value) != int(Len) {
		return TLV{}, errors.New("длинна не совпадает со значением")
	}

	return tlv, nil
}
