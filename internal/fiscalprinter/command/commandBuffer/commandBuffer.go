package commandBuffer

type CommandBufferer interface {
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

func NewCommandBuffer() CommandBufferer {
	return &CommandBuffer{}
}

type CommandBuffer struct{}

func (c CommandBuffer) WriteDate() error {
	panic("implement me")
}

func (c CommandBuffer) WriteTime() error {
	panic("implement me")
}

func (c CommandBuffer) WriteShort() error {
	panic("implement me")
}

func (c CommandBuffer) WriteInt() error {
	panic("implement me")
}

func (c CommandBuffer) WriteLong() error {
	panic("implement me")
}

func (c CommandBuffer) WriteBytes() error {
	panic("implement me")
}

func (c CommandBuffer) WriteByte(b byte) error {
	panic("implement me")
}

func (c CommandBuffer) WriteBoolean() error {
	panic("implement me")
}

func (c CommandBuffer) WriteString() error {
	panic("implement me")
}
