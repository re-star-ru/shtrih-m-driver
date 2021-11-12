package consts

// Типы операций в байтовом виде.
const (
	Income       byte = 1
	ReturnIncome byte = 2
)

// Применять только если целое число! шт или рубль.
const (
	Milligram uint64 = 1000000
	Penny     uint64 = 100 // Минимальная денежная единица(МДЕ)
)

// Налоговые ставки.
const (
	VAT0  byte = 3
	NoVAT byte = 4
)
