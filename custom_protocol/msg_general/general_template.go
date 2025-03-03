package main

var templateStr = `
// Generate By Msg-Gen Tools. By xck.
// DON'T EDIT THIS FILE

package {{.PackageName}}

import (
	"go_study/custom_protocol/process"
	"go_study/custom_protocol/util"
	"bytes"
	"encoding/binary"
)

func (m *{{.MessageName}}) NewIMessage() process.IMessage {
	return new({{.MessageName}})
}

func (m *{{.MessageName}}) Marshal() []byte {
	buf := bytes.NewBuffer([]byte{})
	{{.write_to_buf_txt}}
	return buf.Bytes()
}

func (m *{{.MessageName}}) Unmarshal(buf []byte) {
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

func (m *{{.MessageName}}) SetField(tag uint16, buf []byte) {
	switch tag {
		{{.set_filed_value_txt}}
	}
}

func (m *{{.MessageName}}) GetMessageId() uint16 {
	return 1
}
`

var numberWriteToBufferTemplateStr = `
	if m.{{.FiledName}} > 0 {
		fieldValue := {{.ToByteFunc}}(m.{{.FiledName}})
		// 写入字段tag
		buf.Write(util.Uint16ToByte({{.FieldTag}}))
		// 写入字段长度
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}
`

var stringWriteToBufferTemplateStr = `
	if len(m.{{.FiledName}}) > 0 {
		fieldValue := {{.ToByteFunc}}(m.{{.FiledName}})
		buf.Write(util.Uint16ToByte({{.FieldTag}}))
		buf.Write(util.Uint16ToByte(uint16(len(fieldValue))))
		buf.Write(fieldValue)
	}
`

var setFieldTemplateStr = `
	case {{.FieldTag}}:
		m.{{.FiledName}} = {{.ByteToFieldValueFunc}}(buf)
`
