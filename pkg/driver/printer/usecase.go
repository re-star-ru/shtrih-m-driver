package printer

import "github.com/fess932/shtrih-m-driver/pkg/driver/models"

type Usecase interface {
	OpenShift()  // Открыть смену
	CloseShift() // Закрыть смену

	AddOperationToCheck(op models.Operation)
	CloseCheck()
}

//SellOperationV2(op models.Operation)
