package command

type PrinterCommandImpl struct {
	DefaultTimeout     int
	timeout            int
	resultCode         int
	charsetName        string
	repeatEnabled      bool
	repeadNeeded       bool
	errorReportEnabled bool

	txData []byte
	rxData []byte
}

func (p *PrinterCommandImpl) encodeData() {

}

func newPrinterCommand() *PrinterCommandImpl {
	return &PrinterCommandImpl{
		10000,
		0,
		0,
		"Cp1251", // подумать везде сделать utf-8
		false,
		false,
		true,
		[]byte{},
		[]byte{},
	}
}
