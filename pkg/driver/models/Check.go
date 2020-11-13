package models

// Пакет чека
type CheckPackage struct {
	Operations []Operation `json:"operations"`  // список операций в чеке
	Cash       int64       `json:"cash"`        // сумма оплаты наличными
	Casheless  int64       `json:"casheless"`   // сумма оплаты безналичными
	TaxSystem  byte        `json:"tax_system"`  // система налогообложения
	BottomLine string      `json:"bottom_line"` // нижняя часть чека
}

// Операции в чеке
type Operation struct {
	Type    byte   `json:"type"`    // тип операции
	Amount  int64  `json:"amount"`  // количество товара
	Price   int64  `json:"price"`   // цена в копейках
	Sum     int64  `json:"sum"`     // 	сумма товар * цену
	Subject byte   `json:"subject"` // Предмет рассчета
	Name    string `json:"name"`    // Наименование продукта
}
