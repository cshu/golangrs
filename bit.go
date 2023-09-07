package golangrs

import (
// "math/big"
)

func CopyUint64ToBytes(iBuf uint64, dataBitCount int, outSli []byte) (leftBitCount int, outLen int) {
	leftBitCount = dataBitCount % 8
	outLen = dataBitCount / 8
	iBuf >>= leftBitCount //remove bits that will not be copied (if leftBitCount==0 then it is noop)
	for idx := outLen; idx != 0; {
		idx--
		outSli[idx] = byte(iBuf & 0xff)
		iBuf >>= 8
	}
	return
}

func CopyUint64ToBytesWithPadding(iBuf uint64, dataBitCount int, outSli []byte) int {
	leftBitCount, outLen := CopyUint64ToBytes(iBuf, dataBitCount, outSli)
	if 0 != leftBitCount {
		iBuf <<= (8 - leftBitCount)
		outSli[outLen] = byte(iBuf & 0xff)
		outLen++
	}
	return outLen
}
