package msgpack

import (
	"io"
)

const (
	BIN16 = 0xc5
	UINT8  = 0xcc
	UINT16 = 0xcd
	UINT32 = 0xce
	UINT64 = 0xcf
	STR16   = 0xda
	ARRAY16 = 0xdc

	LEN_INT32 = 4
	LEN_INT64 = 8

	MAX16BIT = 2 << (16 - 1)

	REGULAR_UINT7_MAX  = 2 << (7 - 1)
	REGULAR_UINT8_MAX  = 2 << (8 - 1)
	REGULAR_UINT16_MAX = 2 << (16 - 1)
	REGULAR_UINT32_MAX = 2 << (32 - 1)

	SPECIAL_INT8  = 32
	SPECIAL_INT16 = 2 << (8 - 2)
	SPECIAL_INT32 = 2 << (16 - 2)
	SPECIAL_INT64 = 2 << (32 - 2)
)

type Bytes []byte

// Packs a given value and writes it into the specified writer.
func PackUint8(writer io.Writer, value uint8) (n int, err error) {
	return writer.Write(Bytes{UINT8, value})
}

// Packs a given value and writes it into the specified writer.
func PackUint16(writer io.Writer, value uint16) (n int, err error) {
	return writer.Write(Bytes{UINT16, byte(value >> 8), byte(value)})
}

func PackUint32(writer io.Writer, value uint32) (n int, err error) {
	return writer.Write(Bytes{UINT32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func PackUint64(writer io.Writer, value uint64) (n int, err error) {
	return writer.Write(Bytes{UINT64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func PackBin16(writer io.Writer, value []byte) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{BIN16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write(value)
	return n1 + n2, err
}

func PackStr16(writer io.Writer, value string) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{STR16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write([]byte(value))
	return n1 + n2, err
}

func PackArraySize(writer io.Writer, length uint16) (n int, err error) {
	n, err = writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
	if err != nil {
		return n, err
	}
	return n, nil
}
