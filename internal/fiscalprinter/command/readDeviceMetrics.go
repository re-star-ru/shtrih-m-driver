package command

type ReadDeviceMetrics struct {
	*PrinterCommand
}

func NewReadDeviceMetrics() *ReadDeviceMetrics {
	rc := &ReadDeviceMetrics{NewPrinterCommand()}
	rc.text = "read device metrics"
	rc.commandCode = 252
	return rc
}
