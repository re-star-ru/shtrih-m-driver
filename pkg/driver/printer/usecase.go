package printer

import (
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
)

type Usecase interface {
	OpenShift()            // Открыть смену
	CloseShift()           // Закрыть смену
	ReadShortStatus() byte // Прочитать короткий статус, получить статус

	AddOperationToCheck(op models.Operation) // Добавть операцию в чек
	CloseCheck(chk models.CheckPackage)      // Закрыть чек
}

//SellOperationV2(op models.Operation)
