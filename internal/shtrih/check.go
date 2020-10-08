package shtrih

import (
	"errors"

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
