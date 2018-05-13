package msgpack

import (
	"io"
	"fmt"
)

type (
	Bytes1 [1]byte
	Bytes2 [2]byte
	Bytes4 [4]byte
	Bytes8 [8]byte
)

const (
	NEGFIXNUM     = 0xe0
	FIXMAPMAX     = 0x8f
	FIXARRAYMAX   = 0x9f
	FIXRAWMAX     = 0xbf
	FIRSTBYTEMASK = 0xf
)

func readByte(reader io.Reader) (v uint8, err error) {
	var data Bytes1
	_, e := reader.Read(data[0:])
	if e != nil {
		return 0, e
	}
	return data[0], nil
}

func UnpackUint8(reader io.Reader) (v uint8, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT8  {
		v, err = readByte(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint16(reader io.Reader) (v uint16, n int, err error) {
	var data Bytes2
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint16(data[0]) << 8) | uint16(data[1]), n, nil
}

func UnpackUint16(reader io.Reader) (v uint16, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT16  {
		v, _, err = readUint16(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint32(reader io.Reader) (v uint32, n int, err error) {
	var data Bytes4
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3]), n, nil
}

func UnpackUint32(reader io.Reader) (v uint32, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT32  {
		v, _, err = readUint32(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint64(reader io.Reader) (v uint64, n int, err error) {
	var data Bytes8
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | (uint64(data[3]) << 32) | (uint64(data[4]) << 24) | (uint64(data[5]) << 16) | (uint64(data[6]) << 8) | uint64(data[7]), n, nil
}

func UnpackUint64(reader io.Reader) (v uint64, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT64  {
		v, _, err = readUint64(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func UnpackArraySize(reader io.Reader) (size uint16, err error) {
	c, e := readByte(reader)
	if e != nil {
		return 0, e
	}

	header := uint16(c)
	if header != ARRAY16 {
		return 0, fmt.Errorf("Not Array 16")
	}

	size, _, e = readUint16(reader)
	if e != nil {
		return 0, e
	}

	return size, nil
}

func UnpackStr16(reader io.Reader) (string, error) {
	c, e := readByte(reader)
	if e == nil && c == STR16  {
		size, _, e := readUint16(reader)
		if e == nil {
			value := make([]byte, size)
			n, e := reader.Read(value)
			if e == nil && uint16(n) == size {
				return string(value), nil
			}
		}
	}

	return "", e
}

func UnpackBin16(reader io.Reader) ([]byte, error) {
	c, e := readByte(reader)
	if e == nil && c == BIN16  {
		size, _, e := readUint16(reader)
		if e == nil {
			value := make([]byte, size)
			n, e := reader.Read(value)
			if e == nil && uint16(n) == size {
				return value, nil
			}
		}
	}

	return []byte{}, e
}

