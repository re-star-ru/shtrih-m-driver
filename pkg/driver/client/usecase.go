package client

import "github.com/fess932/shtrih-m-driver/pkg/driver/models"

type Usecase interface {
	Send(frame []byte, cmdLen int) (*models.Frame, error)
	//OperationV2(op models.Operation, password uint32)
	//CloseCheckV2(chk models.CheckPackage, password uint32)
}
