package models

// Пакет чека

type CheckPackage struct {
	CashierINN string      `json:"cashierINN"`
	Operations []Operation `json:"operations"` // Список операций в чеке
	Cash       int64       `json:"cash"`       // Сумма оплаты наличными
	Casheless  int64       `json:"casheless"`  // Сумма оплаты безналичными
	BottomLine string      `json:"bottomLine"` // Нижняя часть чека
	Rounding   byte        `json:"rounding"`   // Округление до рубля, макс 99 копеек
	TaxSystem  byte        `json:"taxSystem"`  // Система налогообложения
	Electronic bool        `json:"electronic"` // Не печатать чек на бумаге
}

// Operation Операции в чеке
type Operation struct {
	Type    byte   `json:"type"`    // Тип операции
	Subject byte   `json:"subject"` // Предмет рассчета
	Amount  int64  `json:"amount"`  // Количество товара
	Price   int64  `json:"price"`   // Цена в копейках
	Sum     int64  `json:"sum"`     // сумма товар * цену
	Name    string `json:"name"`    // Наименование продукта
}
