package TLV

//type Writer struct {
//	buf bufio.Writer
//}

// TEXT в TLV кодировка : CP866

type TLVData []byte

// тип    длинна  значение
//11 04 | 10 00 | 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32

func New(tag uint16, dataLen uint16) TLVData {
	var tlv []byte

	return tlv
}
