package fiscalprinter

func NewFrame() *Frame {
	return &Frame{
		2,
	}
}

type Frame struct {
	STX byte
}

func (f *Frame) GetCrc(data []byte) byte {
	crc := byte(len(data))

	for i := 0; i < len(data); i++ {
		crc ^= data[i]
	}

	return crc
}

//func (f *Frame) encode(data []byte) ([]byte, error) {
//	baos := new ByteArrayOutputStream();
//
//	if (data.length > 255) {
//		throw new Exception("Data length exeeds 256 bytes");
//	} else {
//		baos.write(2);
//		baos.write(data.length);
//		baos.write(data, 0, data.length);
//		baos.write(this.getCrc(data));
//		return baos.toByteArray();
//	}
//}
