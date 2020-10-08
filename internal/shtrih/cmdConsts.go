package shtrih

const (
	// команды принтера
	ReadShortStatus            = 16  // короткий статус принтера
	PrintReportWithoutClearing = 64  // Снять отчет без гашения
	PrintZReport               = 65  // z report?
	WriteTable                 = 30  // запись в таблицу
	ReadTable                  = 31  // чтение из таблицы
	PrintReceiptCopy           = 140 // Печать последнего чека

	// команды фискального накопителя
	FnReadStatus = 65281 // чтение статуса фискального накопителя
	FnReadSerial = 65282 // чтение серийного номера фискального накопителя

	FSDayClose = 65347 // Закрытие фискальной смены
	FSDayOpen  = 65291 // Открытие фискальной смены
)

const (
	NUL = 0x00 // null пустой
	SOH = 0x01 // start of heading начало «заголовка»
	STX = 0x02 // start of text начало «текста»
	ENQ = 0x05 // enquire «Прошу подтверждения!»
	ACK = 0x06 // acknowledgement «Подтверждаю!»
	NAK = 0x15 // negative acknowledgment «Не подтверждаю!»
)
