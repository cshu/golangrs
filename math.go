package golangrs

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func AverageOfFloat64Pair(x float64, y float64) float64 {
	return (x + y) / 2.0
}
func SwapFloat64sIfNotInSortedOrder(x *float64, y *float64) {
	if *x > *y {
		*x, *y = *y, *x
	}
}
func Float64AsExactInt(x float64) (int, bool) {
	retval := int(x)
	exact := float64(retval) == x
	return retval, exact
}
func MaxOfI64Pair(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func RelativeChangeOf(op float64, cl float64) float64 {

	var opRat big.Rat
	var clRat big.Rat
	opRat.SetFloat64(op)
	clRat.SetFloat64(cl)
	rc, _ := new(big.Rat).Quo(new(big.Rat).Sub(&clRat, &opRat), &opRat).Float64()
	return rc
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func LinearInterpolationI64F64(x0 int64, y0 float64, x1 int64, y1 float64, x int64) float64 {
	return y0 + ((float64(x-x0) * (y1 - y0)) / float64(x1-x0))
}

func LinearInterpolationI64F64DiffPointAbove(x0 int64, y0 float64, x1 int64, y1 float64, pointx int64, pointy float64) float64 {
	return pointy - LinearInterpolationI64F64(x0, y0, x1, y1, pointx)
}

func LinearInterpolationI64F64DiffPointBelow(x0 int64, y0 float64, x1 int64, y1 float64, pointx int64, pointy float64) float64 {
	return LinearInterpolationI64F64(x0, y0, x1, y1, pointx) - pointy
}

func AbsByteDiff(x, y byte) byte {
	if x > y {
		return x - y
	} else {
		return y - x
	}
}

func AbsInt64Diff(x, y int64) int64 {
	if x > y {
		return x - y
	} else {
		return y - x
	}
}

func IEEE754Verify(str string, blob []byte, printCont bool) {
	f64, err := strconv.ParseFloat(str, 64)
	CheckErr(err)
	u64 := binary.LittleEndian.Uint64(blob)
	if printCont {
		fmt.Println(str, u64)
	}
	if math.Float64bits(f64) != u64 {
		panic(`IEEE754Verify failed`)
	}
	if f64 != math.Float64frombits(u64) { //?what about NaN? any concern?
		panic(`IEEE754Verify failed`)
	}
}
func IEEE754VerifyPairStrAndBlob(inFile string, printCont bool) {
	pFile, err := os.Open(inFile)
	CheckErr(err)
	defer func() {
		err = pFile.Close()
		CheckErr(err)
	}()
	var bytesbuf [8]byte
	strbuf := make([]byte, math.MaxUint16, math.MaxUint16) //var strbuf [math.MaxUint16]byte
	for {
		_, err = io.ReadFull(pFile, bytesbuf[:2])
		if nil != err {
			if err == io.EOF {
				return
			}
			CheckErrWhichIsNotNil(err)
		}
		strLen := binary.LittleEndian.Uint16(bytesbuf[:2])
		_, err = io.ReadFull(pFile, strbuf[:strLen])
		CheckErr(err)
		_, err = io.ReadFull(pFile, bytesbuf[:])
		CheckErr(err)
		IEEE754Verify(string(strbuf[:strLen]), bytesbuf[:], printCont)
	}
}
func IEEE754GenerateRandomPairStrAndBlob(count int, outFile string, multiplier float64, printCont bool) /* *bytes.Buffer */ { //note generated file could be read by programs written in rust/c++ and see if the verification can pass. If it passes, then this means you can safely pass 8 bytes between these programs (e.g. via tcp socket) for double-precision floating-point number. Say, Rust sends f64 as 8 bytes via TCP to a Go program which read it as float64, supposedly it is ensured to be safe.
	rand.Seed(time.Now().UnixNano())
	pFile, err := os.Create(outFile)
	CheckErr(err)
	defer func() {
		err = pFile.Close()
		CheckErr(err)
	}()
	var bytesbuf [8]byte
	//var buf bytes.Buffer
	for count > 0 {
		count--
		f64 := multiplier * rand.ExpFloat64() //rand.NormFloat64()
		f64str := Float64ToStrUtility(f64)
		u64 := math.Float64bits(f64)
		///
		binary.LittleEndian.PutUint16(bytesbuf[:2], uint16(len(f64str)))
		_, err = pFile.Write(bytesbuf[:2])
		CheckErr(err)
		///
		_, err = pFile.WriteString(f64str)
		CheckErr(err)
		///
		//pFile.WriteString("\000")
		///
		binary.LittleEndian.PutUint64(bytesbuf[:], u64)
		_, err = pFile.Write(bytesbuf[:])
		CheckErr(err)
		///
		if printCont {
			fmt.Println(f64, u64) //fmt.Println(f64str, u64)
		}
	}
	//return &buf
}
