package commandBuffer

type commandBuffer interface {
	WriteDate() error
	WriteTime() error
	//WriteTimeWithSeconds() error wtf????
	WriteShort() error
	WriteInt() error
	WriteLong() error
	WriteBytes() error
	WriteByte(b byte) error

	//WriteBytesWithLength() error wtf"???
	WriteBoolean() error
	WriteString() error
	//WriteStringWithLength() error  wft
}
