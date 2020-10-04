package command

type ReadLongStatus struct {
	*PrinterCommand
}

func NewReadLongStatus() *ReadLongStatus {
	rc := &ReadLongStatus{NewPrinterCommand()}
	rc.text = "Get status"
	rc.commandCode = 17

	return rc
}
