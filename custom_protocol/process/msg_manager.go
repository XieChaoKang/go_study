package process

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

type MessageProcessor struct {
	msgInfo map[uint16]IMessage
	lock    sync.Mutex
}

func (m *MessageProcessor) Register(msg IMessage) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if msg == nil {
		return errors.New(fmt.Sprintf("Register a nil IMessage, check it!!!"))
	}
	msgId := msg.GetMessageId()
	_, ok := m.msgInfo[msgId]
	if ok {
		return errors.New(fmt.Sprintf("this msg id: %v is Registered, may be repeated, check it!!!", msgId))
	}
	m.msgInfo[msgId] = msg
	return nil
}

func (m *MessageProcessor) Marshal(msg IMessage) []byte {
	msgId := msg.GetMessageId()
	bytes := msg.Marshal()
	return append(m.msgIdToByte(msgId), bytes...)
}

func (m *MessageProcessor) UnMarshal(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	msgId := m.byteToMsgId(data[:2])
	m.lock.Lock()
	iMsg, ok := m.msgInfo[msgId]
	m.lock.Unlock()
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found this msg id: %v info", msgId))
	}
	message := iMsg.NewIMessage()
	message.Unmarshal(data[2:])
	return message, nil
}

func (m *MessageProcessor) msgIdToByte(msgId uint16) []byte {
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, msgId)
	return id
}

func (m *MessageProcessor) byteToMsgId(data []byte) uint16 {
	var id uint16
	id = binary.BigEndian.Uint16(data)
	return id
}
