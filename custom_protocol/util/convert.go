package util

import (
	"encoding/binary"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func StringToInt(s string) int {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return int(f)
}

func StringFirstChatToUpper(s string) string {
	split := strings.Split(s, "")
	split[0] = strings.ToUpper(split[0])
	return strings.Join(split, "")
}

func Uint16ToByte(data uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, data)
	return buf
}

func ByteToUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

func Uint32ToByte(data uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, data)
	return buf
}

func ByteToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func Uint64ToByte(data uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, data)
	return buf
}

func ByteToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}
