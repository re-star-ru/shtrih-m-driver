package printer

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

// const PrinterTimeout = 9900 // Дефолтный таймаут ожидания ККТ

type Usecase interface {

	// TODO: OpenShift, CloseShift
	OpenShift()  // Открыть смену
	CloseShift() // Закрыть смену

	ReadShortStatus() byte // Прочитать короткий статус, получить статус

	AddOperationToCheck(op models.Operation)            // Добавть операцию в чек
	CloseCheck(chk models.CheckPackage, dontPrint bool) // Закрыть чек

	CancellationOpenedCheck() // Аннулирование открытого чека

	DontPrintOneCheck()               // пропуск печати одного чека
	WriteCashierINN(INN string) error // Запись Инн кассира
}

//SellOperationV2(op models.Operation)

// printer.writeTable(17, 1, 7, "1") таблица 17 ряд 1 поле 7 значение 1 - не печатать один чек
