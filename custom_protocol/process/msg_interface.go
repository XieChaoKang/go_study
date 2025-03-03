package process

type IMessage interface {
	Marshal() []byte
	Unmarshal(buf []byte)
	SetField(tag uint16, buf []byte)
	GetMessageId() uint16
	NewIMessage() IMessage
}
