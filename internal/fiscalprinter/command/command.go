package command

type PrinterCommand interface {
	GetText() string
	EncodeData() []byte
}
