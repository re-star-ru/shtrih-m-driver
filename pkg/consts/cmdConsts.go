package consts

const (
	// команды принтера
	ReadShortStatus            = 16 // короткий статус принтера
	PrintReportWithoutClearing = 64 // Снять отчет без гашения
	PrintZReport               = 65 // z report?
	WriteTable                 = 30 // запись в таблицу

	ReadFieldInfo    = 46  // Чтение инфо о поле
	ReadTable        = 31  // чтение из таблицы
	PrintReceiptCopy = 140 // Печать последнего чека

	ReadTableStruct = 45 // чтение структуры таблицы

	OpenCheck = 141 // Открытие чека

	PrintSale  = 128 // Добавление операции продажи в чек
	CloseCheck = 133 // Закрытие чека

	// команды фискального накопителя
	FnReadStatus = 65281 // чтение статуса фискального накопителя
	FnReadSerial = 65282 // чтение серийного номера фискального накопителя

	WideRequest = 247 // расширенный запрос

	FnWriteTLV = 65292 // Передать произвольную TLV структуру

	FSDayClose = 65347 // Закрытие фискальной смены
	FSDayOpen  = 65291 // Открытие фискальной смены

	OperationV2  = 65350 // Операция V2  означает начало продажи возврата продажи и тп по нвоому стандарту
	CloseCheckV2 = 65349 // Закрытие чека v2

	SendTLVToOp = 65357 // Передать произвольную TLV структуру привязанную к операции
)

const (
	NUL = 0x00 // null пустой
	SOH = 0x01 // start of heading начало «заголовка»
	STX = 0x02 // start of text начало «текста»
	ENQ = 0x05 // enquire «Прошу подтверждения!»
	ACK = 0x06 // acknowledgement «Подтверждаю!»
	NAK = 0x15 // negative acknowledgment «Не подтверждаю!»
)
