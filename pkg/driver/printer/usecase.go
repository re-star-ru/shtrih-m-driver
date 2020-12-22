package printer

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

// const PrinterTimeout = 9900 // Дефолтный таймаут ожидания ККТ

type Usecase interface {
	OpenShift(c models.Cashier) error  // Открыть смену
	CloseShift(c models.Cashier) error // Закрыть смену

	Print(chk models.CheckPackage) error // Печать чека
	ReadShortStatus() byte               // Прочитать короткий статус, получить статус

	FNOpenedDocumentCancel() error // отменить текущий документ в фн

	// need private \|/
	// CancellationOpenedCheck()                      // Аннулирование открытого чека // todo: сделать приватным
	// AddOperationToCheck(op models.Operation) error // Добавть операцию в чек // todo: сделать приватным
	// CloseCheck(chk models.CheckPackage) error      // Закрыть чек // todo: сделать приватным
	// DontPrintOneCheck() // пропуск печати одного чека // todo: сделать приватным
	// WriteCashierINN(INN string) error              // Запись Инн кассира // todo: сделать приватным
}

//SellOperationV2(op models.Operation)

// printer.writeTable(17, 1, 7, "1") таблица 17 ряд 1 поле 7 значение 1 - не печатать один чек
