// Generate By Msg-Gen Tools. By xck.
// DON'T EDIT THIS FILE

package msg

import (
	"bytes"
	"encoding/binary"
	"go_study/custom_protocol/process"
	"go_study/custom_protocol/util"
)

func (m *Message) NewIMessage() process.IMessage {
	return new(Message)
}

func (m *Message) Marshal() []byte {
	buf := bytes.NewBuffer([]byte{})

	if m.Id > 0 {
		fieldValue := util.Uint32ToByte(m.Id)
		// 写入字段tag
		buf.Write(util.Uint16ToByte(1))
		// 写入字段长度
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}

	if len(m.Name) > 0 {
		fieldValue := util.StringToBytes(m.Name)
		buf.Write(util.Uint16ToByte(2))
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}

	if m.Coins > 0 {
		fieldValue := util.Float32ToByte(m.Coins)
		// 写入字段tag
		buf.Write(util.Uint16ToByte(3))
		// 写入字段长度
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}

	if m.Damoind > 0 {
		fieldValue := util.Float64ToByte(m.Damoind)
		// 写入字段tag
		buf.Write(util.Uint16ToByte(4))
		// 写入字段长度
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}

	return buf.Bytes()
}

func (m *Message) Unmarshal(buf []byte) {
	if len(buf) == 0 {
		return
	}
	for i := uint16(0); i < uint16(len(buf)); {
		// tag 也就是字段id 固定占两位
		tag := binary.BigEndian.Uint16(buf[i : i+2])
		// 字段内容长度 也就是偏移量固定也是占两位
		leng := binary.BigEndian.Uint16(buf[i+2 : i+4])
		// 字段内容根据长度读出来
		fieldValue := buf[i+4 : i+4+leng]
		m.SetField(tag, fieldValue)
		// i 往下一个字段起始位置递增
		i = i + 2 + 2 + leng
	}
}

func (m *Message) SetField(tag uint16, buf []byte) {
	switch tag {

	case 1:
		m.Id = util.ByteToUint32(buf)

	case 2:
		m.Name = util.ByteToString(buf)

	case 3:
		m.Coins = util.ByteToFloat32(buf)

	case 4:
		m.Damoind = util.ByteToFloat64(buf)

	}
}

func (m *Message) GetMessageId() uint16 {
	return 1
}
