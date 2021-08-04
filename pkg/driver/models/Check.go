package models

// Пакет чека
type CheckPackage struct {
	CashierINN string      `json:"cashierINN"`
	Operations []Operation `json:"operations"` // список операций в чеке
	Cash       int64       `json:"cash"`       // сумма оплаты наличными
	Casheless  int64       `json:"casheless"`  // сумма оплаты безналичными
	BottomLine string      `json:"bottomLine"` // нижняя часть чека
	Rounding   byte        `json:"rounding"`   // округление до рубля, макс 99 копеек
	TaxSystem  byte        `json:"taxSystem"`  // система налогообложения
	Electronic bool        `json:"electronic"` // не печатать чек на бумаге
}

// Operation Операции в чеке
type Operation struct {
	Type    byte   `json:"type"`    // тип операции
	Subject byte   `json:"subject"` // Предмет рассчета
	Amount  int64  `json:"amount"`  // количество товара
	Price   int64  `json:"price"`   // цена в копейках
	Sum     int64  `json:"sum"`     // сумма товар * цену
	Name    string `json:"name"`    // Наименование продукта
}
