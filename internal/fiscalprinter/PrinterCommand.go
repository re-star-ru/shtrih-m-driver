package fiscalprinter

type PrinterCommand struct {
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

func (p *PrinterCommand) encodeData() {

}

func newPrinterCommand() *PrinterCommand {
	return &PrinterCommand{
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
