package type_convert

import (
	"encoding/binary"
)

var DefaultByteOrder binary.ByteOrder //默认字节序

func bytesToInt8(bs []byte) (i8 int8, err error) {
	if len(bs) < 1 { //空的[]byte 转为　0
		return 0, nil
	}
	return int8(bs[0]), nil
}

func bytesToInt16(bs []byte) (i16 int16, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return int16(getByteAt(bs, l, 0)) | int16(getByteAt(bs, l, 1))<<8, nil
	}
	return int16(getByteAt(bs, l, 1)) | int16(getByteAt(bs, l, 0))<<8, nil
}

func bytesToInt32(bs []byte) (i32 int32, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return int32(getByteAt(bs, l, 0)) | int32(getByteAt(bs, l, 1))<<8 | int32(getByteAt(bs, l, 2))<<16 | int32(getByteAt(bs, l, 3))<<24, nil
	}
	return int32(getByteAt(bs, l, 3)) | int32(getByteAt(bs, l, 2))<<8 | int32(getByteAt(bs, l, 1))<<16 | int32(getByteAt(bs, l, 0))<<24, nil
}

func bytesToInt64(bs []byte) (i64 int64, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return int64(getByteAt(bs, l, 0)) | int64(getByteAt(bs, l, 1))<<8 |
			int64(getByteAt(bs, l, 2))<<16 | int64(getByteAt(bs, l, 3))<<24 |
			int64(getByteAt(bs, l, 4))<<32 | int64(getByteAt(bs, l, 5))<<40 |
			int64(getByteAt(bs, l, 6))<<48 | int64(getByteAt(bs, l, 7))<<56, nil
	}
	return int64(getByteAt(bs, l, 7)) | int64(getByteAt(bs, l, 6))<<8 |
		int64(getByteAt(bs, l, 5))<<16 | int64(getByteAt(bs, l, 4))<<24 |
		int64(getByteAt(bs, l, 3))<<32 | int64(getByteAt(bs, l, 2))<<40 |
		int64(getByteAt(bs, l, 1))<<48 | int64(getByteAt(bs, l, 0))<<56, nil
}

func bytesToUint8(bs []byte) (u8 uint8, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	return bs[0], nil
}

func bytesToUint16(bs []byte) (u16 uint16, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return uint16(getByteAt(bs, l, 0)) | uint16(getByteAt(bs, l, 1))<<8, nil
	}
	return uint16(getByteAt(bs, l, 1)) | uint16(getByteAt(bs, l, 0))<<8, nil
}

func bytesToUint32(bs []byte) (u32 uint32, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return uint32(getByteAt(bs, l, 0)) | uint32(getByteAt(bs, l, 1))<<8 | uint32(getByteAt(bs, l, 2))<<16 | uint32(getByteAt(bs, l, 3))<<24, nil
	}
	return uint32(getByteAt(bs, l, 3)) | uint32(getByteAt(bs, l, 2))<<8 | uint32(getByteAt(bs, l, 1))<<16 | uint32(getByteAt(bs, l, 0))<<24, nil
}

func bytesToUint64(bs []byte) (u64 uint64, err error) {
	l := len(bs)
	if l < 1 {
		return 0, nil
	}
	if DefaultByteOrder == binary.LittleEndian {
		return uint64(getByteAt(bs, l, 0)) | uint64(getByteAt(bs, l, 1))<<8 |
			uint64(getByteAt(bs, l, 2))<<16 | uint64(getByteAt(bs, l, 3))<<24 |
			uint64(getByteAt(bs, l, 4))<<32 | uint64(getByteAt(bs, l, 5))<<40 |
			uint64(getByteAt(bs, l, 6))<<48 | uint64(getByteAt(bs, l, 7))<<56, nil
	}
	return uint64(getByteAt(bs, l, 7)) | uint64(getByteAt(bs, l, 6))<<8 |
		uint64(getByteAt(bs, l, 5))<<16 | uint64(getByteAt(bs, l, 4))<<24 |
		uint64(getByteAt(bs, l, 3))<<32 | uint64(getByteAt(bs, l, 2))<<40 |
		uint64(getByteAt(bs, l, 1))<<48 | uint64(getByteAt(bs, l, 0))<<56, nil
}

func getByteAt(bs []byte, l int, pos int) byte {
	if pos >= l {
		return 0
	}
	return bs[pos]
}
