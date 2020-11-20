package models

// Пакет чека
type CheckPackage struct {
	CashierINN string      `json:"cashier_inn"`
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
	Sum     int64  `json:"sum"`     // сумма товар * цену
	Subject byte   `json:"subject"` // Предмет рассчета
	Name    string `json:"name"`    // Наименование продукта
}

//type CloseCheckModel struct {
//	Payment1      []byte
//	Payment2      []byte
//	Payment3      []byte
//	Payment4      []byte
//	Payment5      []byte
//	Payment6      []byte
//	Payment7      []byte
//	Payment8      []byte
//	Payment9      []byte
//	Payment10     []byte
//	Payment11     []byte
//	Payment12     []byte
//	Payment13     []byte
//	Payment14     []byte
//	Payment15     []byte
//	Payment16     []byte // 80
//	RubleRounding []byte // 81
//	Tax1          []byte
//	Tax2          []byte
//	Tax3          []byte
//	Tax4          []byte
//	Tax5          []byte
//	Tax6          []byte // 30 + 81 = 111
//	TaxSystem     []byte // 112
//	Text          []byte // 176 + passsword(4) = 180 + cmdLen(2) = 182
//}
//
//func (c *CloseCheckModel) Bytes() []byte {
//
//	b := make([]byte, 176)
//
//	copy(b[0:5], c.Payment1)
//	copy(b[5:10], c.Payment2)
//	copy(b[10:15], c.Payment3)
//	copy(b[15:20], c.Payment4)
//	copy(b[20:25], c.Payment5)
//	copy(b[25:30], c.Payment6)
//	copy(b[30:35], c.Payment7)
//	copy(b[35:40], c.Payment8)
//	copy(b[40:45], c.Payment9)
//	copy(b[45:50], c.Payment10)
//	copy(b[50:55], c.Payment11)
//	copy(b[55:60], c.Payment12)
//	copy(b[60:65], c.Payment13)
//	copy(b[65:70], c.Payment14)
//	copy(b[70:75], c.Payment15)
//	copy(b[75:80], c.Payment16)
//
//	return b
//}
//
//func NewCloseCheckModel() *CloseCheckModel {
//	return &CloseCheckModel{
//		Payment1:      make([]byte, 5),
//		Payment2:      make([]byte, 5),
//		Payment3:      make([]byte, 5),
//		Payment4:      make([]byte, 5),
//		Payment5:      make([]byte, 5),
//		Payment6:      make([]byte, 5),
//		Payment7:      nil,
//		Payment8:      nil,
//		Payment9:      nil,
//		Payment10:     nil,
//		Payment11:     nil,
//		Payment12:     nil,
//		Payment13:     nil,
//		Payment14:     nil,
//		Payment15:     nil,
//		Payment16:     nil,
//		RubleRounding: nil,
//		Tax1:          nil,
//		Tax2:          nil,
//		Tax3:          nil,
//		Tax4:          nil,
//		Tax5:          nil,
//		Tax6:          nil,
//		TaxSystem:     nil,
//		Text:          nil,
//	}
//}
