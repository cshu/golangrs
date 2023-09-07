package golangrs

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCopyUint64ToBytes(t *testing.T) {
	var x uint64 = 0b10000001100101011
	var usedBitCount int = 17
	var outSli []byte = make([]byte, 0x100)
	leftBitCount, outLen := CopyUint64ToBytes(x, usedBitCount, outSli)
	if leftBitCount != 1 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if outLen != 2 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b10000001, 0b10010101}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	x <<= 1
	leftBitCount, outLen = CopyUint64ToBytes(x, usedBitCount+1, outSli)
	if leftBitCount != 2 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if outLen != 2 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b10000001, 0b10010101}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	x = 0b1000000100110101100
	leftBitCount, outLen = CopyUint64ToBytes(x, 19, outSli)
	if leftBitCount != 3 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if outLen != 2 {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b10000001, 0b110101}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	z := int64(-1)
	x = uint64(z)
	fmt.Println(`greatest uint64`, x)
	x <<= 2
	fmt.Println(x)
	leftBitCount, outLen = CopyUint64ToBytes(x, 9, outSli)
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b11111110}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	leftBitCount, outLen = CopyUint64ToBytes(x, 17, outSli)
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b11111111, 0b11111110}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	x >>= 2
	fmt.Println(x)
	leftBitCount, outLen = CopyUint64ToBytes(x, 17, outSli)
	if !reflect.DeepEqual(outSli[:outLen], []byte{0b11111111, 0b11111111}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	leftBitCount, outLen = CopyUint64ToBytes(x, 64, outSli)
	if !reflect.DeepEqual(outSli[:2], []byte{0b00111111, 0b11111111}) {
		t.Errorf(`CopyUint64ToBytes failure`)
	}
	x <<= 2
	fmt.Println(x)
	fmt.Println(`CopyUint64ToBytes tested`)
}
