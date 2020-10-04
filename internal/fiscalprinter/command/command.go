package command

type PrinterCommander interface {
	GetCode() int
	GetText() string
	EncodeData() []byte
}
