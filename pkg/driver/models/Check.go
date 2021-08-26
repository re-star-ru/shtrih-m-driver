package models

// CheckPackage это чек для одной кассы.
type CheckPackage struct {
	CashierINN string
	Operations []Operation // Список операций в чеке
	Cash       uint64      // Сумма оплаты наличными
	Digital    uint64      // Сумма оплаты безналичными
	Rounding   byte        // Округление до рубля, макс 99 копеек
	TaxSystem  byte        // Система налогообложения
	NotPrint   bool        // Не печатать чек на бумаге
}

// Operation Операции в чеке.
type Operation struct {
	Type    byte   // Тип операции
	Subject byte   // Предмет рассчета
	Amount  uint64 // Количество товара
	Price   uint64 // Цена в копейках
	Sum     uint64 // сумма товар * цену
	Name    string // Наименование продукта
}
