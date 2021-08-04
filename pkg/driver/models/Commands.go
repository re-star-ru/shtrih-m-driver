package models

import "time"

const (
	// команды принтера
	ReadShortStatus            = 16 // короткий статус принтера
	PrintReportWithoutClearing = 64 // Снять отчет без гашения
	WriteTable                 = 30 // запись в таблицу

	ReadFieldInfo    = 46  // Чтение инфо о поле
	ReadTable        = 31  // чтение из таблицы
	PrintReceiptCopy = 140 // Печать последнего чека

	ReadTableStruct = 45 // чтение структуры таблицы

	PrintSale = 128 // Добавление операции продажи в чек

	CloseCheck = 133 // Закрытие чека

	CancellationCheck = 136 // Отмена текущего чека

	// команды фискального накопителя
	FnReadStatus = 65281 // чтение статуса фискального накопителя

	WideRequest = 247 // расширенный запрос

	StartOpenShift = 65345 // начать открытие смены
	OpenShift      = 224   // открыть смену

	StartCloseShift = 65346 // начало закрытия смены
	ZReport         = 65    // суточный отчет с гашением, (закрытие смены)

	FNCancelCurrentDocument = 65288 // отменить открытый документ в фн

	FNCloseShift = 65347 // Закрытие фискальной смены

	OperationV2  = 65350 // Операция V2  означает начало продажи возврата продажи и тп по нвоому стандарту
	CloseCheckV2 = 65349 // Закрытие чека v2

	FnWriteTLV = 65292 // 0C

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

const (
	DefaultAttemptTimeout = 1000 * time.Millisecond
	MaxENQAttempts        = 3
)
