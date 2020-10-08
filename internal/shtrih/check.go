package shtrih

import (
	"math"
	"math/big"
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
	d big.Float


	BigDecimal Cash = BigDecimal.ZERO;

	/**
	 * Сумма электронной оплаты
	 */
	@Attribute(required = false)
	public BigDecimal ElectronicPayment = BigDecimal.ZERO;

	/**
	 * Сумма предоплатой (зачетом аванса)
	 */
	@Attribute(required = false)
	public BigDecimal AdvancePayment = BigDecimal.ZERO;

	/**
	 * Сумма постоплатой (в кредит)
	 */
	@Attribute(required = false)
	public BigDecimal Credit = BigDecimal.ZERO;

	/**
	 * Сумма встречным предоставлением
	 */
	@Attribute(required = false)
	public BigDecimal CashProvision = BigDecimal.ZERO;
}

type Position struct {
	typeString string
}
