package shtrih

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/fess932/shtrih-m-driver/internal/shtrih/check"

	"golang.org/x/text/encoding/charmap"

	"github.com/shopspring/decimal"
)

type CheckPackage struct {
	CashierName  string // ФИО и должность уполномоченного лица для проведения операции
	CashierVATIN string // ИНН уполномоченного лица для проведения операции

	/**
	 * Тип расчета
	 * 1 - Приход
	 * 2 - Возврат прихода
	 * 3 - Расход
	 * 4 - Возврат расхода
	 */
	PaymentType int

	/**
	 * Код системы налогообложения
	 * 0 Общая
	 * 1 Упрощенная Доход
	 * 2 Упрощенная Доход минус Расход
	 * 3 Единый налог на вмененный доход
	 * 4 Единый сельскохозяйственный налог
	 * 5 Патентная система налогообложения
	 */
	TaxVariant int

	/**
	 * Email покупателя
	 */
	CustomerEmail string

	/**
	 * Телефонный номер покупателя
	 */
	CustomerPhone string

	/**
	 * Адрес электронной почты отправителя чека
	 */
	SenderEmail string

	/**
	 * Адрес проведения расчетов
	 */
	AddressSettle string

	/**
	 * Место проведения расчетов
	 */
	PlaceSettle string

	/**
	 * Позиции в новом чеке
	 */
	Positions []Position

	/**
	 * Параметры закрытия чека. Чек коррекции может быть оплачен только одним видом оплаты
	 * и без сдачи.
	 */

	/**
	 * Сумма наличной оплаты
	 */
	Cash decimal.Decimal

	/**
	 * Сумма электронной оплаты
	 */
	ElectronicPayment decimal.Decimal

	/**
	* Сумма предоплатой (зачетом аванса)
	 */
	AdvancePayment decimal.Decimal

	/**
	* Сумма постоплатой (в кредит)
	 */
	Credit decimal.Decimal
	///**
	// * Сумма встречным предоставлением
	// */
	//@Attribute(required = false)
	//public BigDecimal CashProvision = BigDecimal.ZERO;
}

type Position struct {
	typeString string
}

type FiscalString struct {
	Name              string          // Наименование товара
	Quantity          decimal.Decimal // Количество товара
	PriceWithDiscount decimal.Decimal // Цена единицы товара с учетом скидок/наценок
	SumWithDiscount   decimal.Decimal // 	 * Конечная сумма по позиции чека с учетом всех скидок/наценок

	/**
	 * Ставка НДС. Список значений:
	 *  "none" - БЕЗ НДС
	 *  "20" - НДС 20
	 *  "10" - НДС 10
	 *  "0" - НДС 0
	 *  "10/110" - расч. ставка 10/110
	 *  "20/120" - расч. ставка 20/120
	 */
	Tax string

	SignMethodCalculation int    // Признак способа расчета
	SignCalculationObject int    // Признак предмета расчета
	MeasurementUnit       string // Единица измерения предмета расчета

	//GoodCodeData GoodCodeData // Данные кода товарной номенклатуры

}

func (f *FiscalString) getTax() (int, error) {
	switch f.Tax {
	case "20":
		return 1, nil
	case "10":
		return 2, nil
	case "20/120":
		return 3, nil
	case "10/110":
		return 4, nil
	case "0":
		return 5, nil
	case "none":
		return 6, nil
	default:
		return 0, errors.New("Неизвестный тип налоговой ставки: " + f.Tax)
	}
}

type TextString struct {
	Text       string // Строка с произвольным текстом
	FontNumber int    // Строка с произвольным текстом
}

type BarcodeString struct {
	BarcodeType string // Строка, определяющая тип штрихкода
	Barcode     string // Значение штрихкода
}

////////////////////////////////////// Продажа

func (p *Printer) SellOperationV2() {
	data, cmdLen := p.createCommandData(OperationV2)
	buf := bytes.NewBuffer(data)
	// Запись типа операции
	buf.WriteByte(check.Income) // Тип операции

	// Запись количества товара
	// Количество записывается в миллиграммах
	amount, err := intToBytesWithLen(2*check.Milligram, 6)
	if err != nil {
		p.logger.Fatal(err)
	}
	p.logger.Debug("amount:\n", hex.Dump(amount))
	buf.Write(amount)

	// запись цены товара
	// цена записывается в копейках
	price, err := intToBytesWithLen(1, 5) // одна копейка
	if err != nil {
		p.logger.Fatal(err)
	}
	buf.Write(price)

	// запись суммы товара
	// Сумма записывается в копейках
	summ, err := intToBytesWithLen(2, 5) // две копейки
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

	const ff = 0xff
	buf.Write([]byte{ff, ff, ff, ff, ff}) // если нет налога надо отправлять 0xff
	// Запись налоговой ставки
	buf.WriteByte(check.VAT0)
	// Запись номера отдела
	buf.WriteByte(1)

	// Запись признака способа рассчета
	buf.WriteByte(check.FullPayment)

	// Запись признака предмета рассчета
	buf.WriteByte(check.Service)

	// Запись название товара 0 - 128 байт строка
	str, err := charmap.Windows1251.NewEncoder().String("Товар 1 charmap1251")
	//str := "Товар 1"
	//charmap.Windows1251.NewEncoder().String()
	if err != nil {
		p.logger.Fatal(err)
	}
	rStrBytes := make([]byte, 128)
	copy(rStrBytes, []byte(str))

	buf.Write(rStrBytes[:128])

	p.logger.Debug("длинна сообщения в байтах: ", buf.Len())
	p.logger.Debug("\n", hex.Dump(buf.Bytes()))

	p.logger.Debug("cmdlen", cmdLen)
	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.bytes()))

}

///////////////////////////////////// Закрытие чека

func (p *Printer) CloseCheckV2() {
	data, cmdLen := p.createCommandData(CloseCheckV2)
	buf := bytes.NewBuffer(data)
	p.logger.Debug("cmdlen:", cmdLen)

	// запись суммы наличных
	cash, err := intToBytesWithLen(2, 5)
	if err != nil {
		p.logger.Fatal(err)
	}
	buf.Write(cash)

	// запись суммы типа оплаты 2 - безнал
	casheless, err := intToBytesWithLen(0, 5)
	if err != nil {
		p.logger.Fatal(err)
	}

	for i := 1; i < 16; i++ {
		buf.Write(casheless) // 2 - 16
	}

	buf.WriteByte(0) // округление до рубля

	buf.Write(casheless) // налог 1
	buf.Write(casheless) // налог 2
	buf.Write(casheless) // налог 3
	buf.Write(casheless) // налог 4
	buf.Write(casheless) // налог 5
	buf.Write(casheless) // налог 6

	buf.WriteByte(check.ENVD) //система налогоообложения, биты а не байт

	// Запись название товара 0 - 128 байт строка
	str, err := charmap.Windows1251.NewEncoder().String("нижняя часть чека 64 байта")
	//str := "Товар 1"
	//charmap.Windows1251.NewEncoder().String()
	if err != nil {
		p.logger.Fatal(err)
	}
	rStrBytes := make([]byte, 64)
	copy(rStrBytes, []byte(str))

	buf.Write(rStrBytes[:64])

	p.logger.Debug("len: ", buf.Len())

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.bytes()))
}

func intToBytesWithLen(val int64, bytesLen int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}

	return buf.Bytes()[:bytesLen], nil
}
