package deprecated

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/consts"

	"golang.org/x/text/encoding/charmap"
)

//type CheckPackage struct {
//	Operations []Operation `json:"operations"`  // список операций в чеке
//	Cash       int64       `json:"cash"`        // сумма оплаты наличными
//	Casheless  int64       `json:"casheless"`   // сумма оплаты безналичными
//	TaxSystem  byte        `json:"tax_system"`  // система налогообложения
//	BottomLine string      `json:"bottom_line"` // нижняя часть чека
//}
//
//type Operation struct {
//	Type    byte   `json:"type"`    // тип операции
//	Amount  int64  `json:"amount"`  // количество товара
//	Price   int64  `json:"price"`   // цена в копейках
//	Sum     int64  `json:"sum"`     // 	сумма товар * цену
//	Subject byte   `json:"subject"` // Предмет рассчета
//	Name    string `json:"name"`    // Наименование продукта
//}

////////////////////////////////////// Операция v2

func (p *Printer) SellOperationV2(op models.Operation) {
	data, cmdLen := p.createCommandData(models.OperationV2)
	buf := bytes.NewBuffer(data)

	// Запись типа операции
	buf.WriteByte(op.Type)

	// Запись количества товара
	// Количество записывается в миллиграммах
	amount, err := intToBytesWithLen(op.Amount*consts.Milligram, 6)
	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug("amount:\n", hex.Dump(amount))
	buf.Write(amount)

	// запись цены товара
	// цена записывается в копейках
	price, err := intToBytesWithLen(op.Price, 5) // одна копейка
	if err != nil {
		p.logger.Fatal(err)
	}
	buf.Write(price)

	// запись суммы товара
	// Сумма записывается в копейках
	summ, err := intToBytesWithLen(op.Sum, 5) // две копейки
	if err != nil {
		p.logger.Fatal(err)
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
		p.logger.Fatal(err)
	}
	rStrBytes := make([]byte, 128)
	copy(rStrBytes, str)

	buf.Write(rStrBytes[:128])

	p.logger.Debug("длинна сообщения в байтах: ", buf.Len())
	p.logger.Debug("\n", hex.Dump(buf.Bytes()))

	p.logger.Debug("cmdlen", cmdLen)
	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

}

///////////////////////////////////// Закрытие чека

func (p *Printer) CloseCheckV2(chk models.CheckPackage) {
	data, cmdLen := p.createCommandData(models.CloseCheckV2)
	buf := bytes.NewBuffer(data)
	p.logger.Debug("cmdlen:", cmdLen)

	// запись суммы наличных - типа оплаты 1
	cash, err := intToBytesWithLen(chk.Cash, 5)
	if err != nil {
		p.logger.Fatal(err)
	}
	buf.Write(cash)

	// запись суммы типа оплаты 2 - безнал
	casheless, err := intToBytesWithLen(chk.Casheless, 5)
	if err != nil {
		p.logger.Fatal(err)
	}
	buf.Write(casheless)

	for i := 2; i < 16; i++ {
		buf.Write(make([]byte, 5)) // 3 - 16
	}

	buf.WriteByte(0) // округление до рубля

	for i := 0; i < 5; i++ {
		buf.Write(make([]byte, 5)) // налог 1-6

	}
	//buf.Write(casheless) // налог 1
	//buf.Write(casheless) // налог 2
	//buf.Write(casheless) // налог 3
	//buf.Write(casheless) // налог 4
	//buf.Write(casheless) // налог 5
	//buf.Write(casheless) // налог 6

	buf.WriteByte(chk.TaxSystem) //система налогоообложения, биты а не байт

	// Запись название товара 0 - 128 байт строка
	str, err := charmap.Windows1251.NewEncoder().String(chk.BottomLine)
	if err != nil {
		p.logger.Fatal(err)
	}
	rStrBytes := make([]byte, 64)
	copy(rStrBytes, str)

	buf.Write(rStrBytes[:64])

	p.logger.Debug("len: ", buf.Len())

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))
}

func intToBytesWithLen(val int64, bytesLen int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}

	return buf.Bytes()[:bytesLen], nil
}
