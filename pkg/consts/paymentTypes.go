package consts

// Признак способа рассчета
const (
	PrePayment100 byte = 1
	PrePayment    byte = 2
	Imprest       byte = 3
	FullPayment   byte = 4 // предоплата
)

// Признак предмета рассчета/ subject
const (
	Goods       byte = 1 // товар
	ExciseGoods byte = 2 // акцизный товар
	Job         byte = 3 // работа
	Service     byte = 4 // услуга
)
