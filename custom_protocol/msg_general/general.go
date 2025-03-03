package main

import (
	"awesomeProject1/general_"
	"awesomeProject1/util"
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type msgFieldType string

const (
	msgFieldTypeUint32  msgFieldType = "uint32"
	msgFieldTypeFloat32 msgFieldType = "float32"
	msgFieldTypeFloat64 msgFieldType = "float64"
	msgFieldTypeString  msgFieldType = "string"
)

type MsgTagInfo struct {
	fieldType msgFieldType
	fieldTag  int
	filedName string // 注意 这个必须遵循驼峰命名规则
}

func (m *MsgTagInfo) generalFieldToByteFuncTxt() string {
	txt := ""
	switch m.fieldType {
	case msgFieldTypeUint32:
		txt = "util.Uint32ToByte"
	case msgFieldTypeFloat32:
		txt = "util.Float32ToByte"
	case msgFieldTypeFloat64:
		txt = "util.Float64ToByte"
	case msgFieldTypeString:
		txt = "util.StringToBytes"
	}
	return txt
}

func (m *MsgTagInfo) generalByteToFieldValueFuncTxt() string {
	txt := ""
	switch m.fieldType {
	case msgFieldTypeUint32:
		txt = "util.ByteToUint32"
	case msgFieldTypeFloat32:
		txt = "util.ByteToFloat32"
	case msgFieldTypeFloat64:
		txt = "util.ByteToFloat64"
	case msgFieldTypeString:
		txt = "util.ByteToString"
	}
	return txt
}

func (m *MsgTagInfo) convertPrimitiveFieldName() string {
	// 这里是按照驼峰命名的 所以按照驼峰形式解析转换即可
	// 按照下划线切割 同时把每个词的单词换成首字母即可还原回原本的字段名称
	split := strings.Split(m.filedName, "_")
	var name []string
	for _, char := range split {
		name = append(name, util.StringFirstChatToUpper(char))
	}
	return strings.Join(name, "")
}

func (m *MsgTagInfo) generalFieldWriteToBufferTxt() string {
	primitiveFiledName := m.convertPrimitiveFieldName()
	filedToByteFuncTxt := m.generalFieldToByteFuncTxt()
	template := ``
	switch m.fieldType {
	case msgFieldTypeUint32, msgFieldTypeFloat32, msgFieldTypeFloat64:
		template = strings.Clone(numberWriteToBufferTemplateStr)
	case msgFieldTypeString:
		template = strings.Clone(stringWriteToBufferTemplateStr)
	}
	template = strings.ReplaceAll(template, "{{.FiledName}}", primitiveFiledName)
	template = strings.ReplaceAll(template, "{{.ToByteFunc}}", filedToByteFuncTxt)
	template = strings.ReplaceAll(template, "{{.FieldTag}}", fmt.Sprintf("%v", m.fieldTag))
	return template
}

func (m *MsgTagInfo) generalSetFieldValueTxt() string {
	primitiveFiledName := m.convertPrimitiveFieldName()
	byteToFieldValueFuncTxt := m.generalByteToFieldValueFuncTxt()
	template := strings.Clone(setFieldTemplateStr)
	template = strings.ReplaceAll(template, "{{.FiledName}}", primitiveFiledName)
	template = strings.ReplaceAll(template, "{{.ByteToFieldValueFunc}}", byteToFieldValueFuncTxt)
	template = strings.ReplaceAll(template, "{{.FieldTag}}", fmt.Sprintf("%v", m.fieldTag))
	return template
}

func GetMsgTag(msg general_.Message) []string {
	t := reflect.ValueOf(msg).Type()
	var msgTags []string
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		structTag := structField.Tag
		msgTag := structTag.Get("msg_tag")
		println(msgTag)
		msgTags = append(msgTags, msgTag)
	}
	return msgTags
}

func passMsgTag(tag string) (*MsgTagInfo, error) {
	split := strings.Split(tag, ",")
	if len(split) != 3 {
		return nil, errors.New(fmt.Sprintf("msg tag len is illegal, check it!! tag txt: %v", tag))
	}
	info := &MsgTagInfo{
		msgFieldType(split[0]),
		util.StringToInt(split[1]),
		split[2],
	}
	return info, nil
}

func generalWriteToBufferTxt(tagInfos []*MsgTagInfo) string {
	var writeToBufTxt []string
	for _, tagInfo := range tagInfos {
		writeToBufTxt = append(writeToBufTxt, tagInfo.generalFieldWriteToBufferTxt())
	}
	join := strings.Join(writeToBufTxt, "\n")
	return join
}

func generalSetFiledValueTxt(tagInfos []*MsgTagInfo) string {
	var setFiledValueTxt []string
	for _, tagInfo := range tagInfos {
		setFiledValueTxt = append(setFiledValueTxt, tagInfo.generalSetFieldValueTxt())
	}
	join := strings.Join(setFiledValueTxt, "")
	return join
}

func general(msgTag []string, msgName string) string {
	var tagInfos []*MsgTagInfo
	for _, tag := range msgTag {
		tagInfo, err := passMsgTag(tag)
		if err != nil {
			panic(fmt.Sprintf("General msg: %v pass mag tag err: %v", msgName, err.Error()))
		}
		tagInfos = append(tagInfos, tagInfo)
	}
	tempTemplateStr := strings.Clone(templateStr)
	tempTemplateStr = strings.ReplaceAll(tempTemplateStr, "{{.MessageName}}", msgName)
	tempTemplateStr = strings.ReplaceAll(tempTemplateStr, "{{.write_to_buf_txt}}", generalWriteToBufferTxt(tagInfos))
	tempTemplateStr = strings.ReplaceAll(tempTemplateStr, "{{.set_filed_value_txt}}", generalSetFiledValueTxt(tagInfos))
	return tempTemplateStr
}

var removeChat = '"'

func getMsgTagFromFieldTag(tag string) string {
	if len(tag) == 0 {
		return ""
	}
	tag = strings.ReplaceAll(tag, "`", "")
	tags := strings.Split(tag, " ")
	for _, tempTag := range tags {
		tagInfo := strings.Split(tempTag, ":")
		if len(tagInfo) != 2 {
			continue
		}
		if tagInfo[0] == "msg_tag" {
			return strings.ReplaceAll(tagInfo[1], string(removeChat), "")
		}
	}
	return ""
}

func General(filePath string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("General filePath: %v ParseFile err: %v", filePath, err.Error()))
		return
	}
	if len(f.Scope.Objects) == 0 {
		println("not found obj")
		return
	}
	ext := filepath.Ext(filePath)
	newFilePath, _ := filepath.Abs(fmt.Sprintf("%s_auto_gen%s", filePath[:len(filePath)-len(ext)], ext))
	file, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("General open file: %v err: %v", newFilePath, err.Error()))
	}
	defer file.Close()

	var generalTxt []string
	for objName, object := range f.Scope.Objects {
		modelTypeSpec, ok := object.Decl.(*ast.TypeSpec)
		if !ok {
			panic("cannot find model type")
			return
		}
		modelStructType, ok := modelTypeSpec.Type.(*ast.StructType)
		if !ok {
			panic("cannot translate model type to struct type")
			return
		}
		var msgTags []string
		for _, field := range modelStructType.Fields.List {
			if len(field.Names) > 0 {

			}
			msgTagFromFieldTag := getMsgTagFromFieldTag(field.Tag.Value)
			if len(msgTagFromFieldTag) == 0 {
				continue
			}
			msgTags = append(msgTags, msgTagFromFieldTag)
		}
		s := general(msgTags, objName)
		generalTxt = append(generalTxt, s)
	}
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	_, err = write.WriteString(strings.Join(generalTxt, "\n"))
	if err != nil {
		panic(fmt.Sprintf("General  WriteString err: %v", err.Error()))
	}
	//Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		panic(fmt.Sprintf("General Flush err: %v", err.Error()))
		return
	}
}
