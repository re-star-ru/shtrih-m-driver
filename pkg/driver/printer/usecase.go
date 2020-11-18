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

	AddOperationToCheck(op models.Operation) // Добавть операцию в чек
	CloseCheck(chk models.CheckPackage)      // Закрыть чек

	CancellationOpenedCheck() // Аннулирование открытого чека
}

//SellOperationV2(op models.Operation)
