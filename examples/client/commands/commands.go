package commands

// KKT Commands
const (
	ShortStatus byte = 0x10
	ZReport     byte = 0x41
)

// Fn Commands Start with FF
const (
	FnOperationV2      byte = 0x46
	FnCloseCheckV2     byte = 0x45
	FnCancelFNDocument byte = 0x08

	FnBeginOpenSession byte = 0x41 // start then
	FnWriteTLV         byte = 0x0C // send tlv then
	FnOpenSession      byte = 0x0B // end open
)

var defaultPassword = []byte{0x1E, 0x00, 0x00, 0x00}

func CreateShortStatus() (cmdID byte, cmdData []byte) {
	return ShortStatus, defaultPassword
}
