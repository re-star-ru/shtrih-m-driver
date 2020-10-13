package check

// Типы чеков, FiscalReceiptType
const (
	SMFPTR_RT_SALE    = 100 // приход
	SMFPTR_RT_BUY     = 101 // расход
	SMFPTR_RT_RETSALE = 102 // возврат прихода
	SMFPTR_RT_RETBUY  = 103 // возврат расхода
)

type Check struct {
	CashierName       string // Имя кассира
	FiscalReceiptType int    // Тип чека
	TaxVariant        int    // Система налогообложения
}

func New() *Check {
	return &Check{}
}
