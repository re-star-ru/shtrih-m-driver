package models

// CheckPackage это чек для одной кассы.
type CheckPackage struct {
	CashierINN string
	Operations []Operation // Список операций в чеке
	Cash       int64       // Сумма оплаты наличными
	Digital    int64       // Сумма оплаты безналичными
	Rounding   byte        // Округление до рубля, макс 99 копеек
	TaxSystem  byte        // Система налогообложения
	NotPrint   bool        // Не печатать чек на бумаге
}

// Operation Операции в чеке.
type Operation struct {
	Type    byte   // Тип операции
	Subject byte   // Предмет рассчета
	Amount  int64  // Количество товара
	Price   int64  // Цена в копейках
	Sum     int64  // сумма товар * цену
	Name    string // Наименование продукта
}
