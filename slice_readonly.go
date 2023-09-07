package golangrs

import (
	"math"
	"sort"
)

func DoesSliceContain[E comparable](sli []E, v E) bool {
	for _, elem := range sli {
		if elem == v {
			return true
		}
	}
	return false
}

func IndexOfElemInSliceWithPredicate[E any](sli []E, cb func(E) bool) int {
	for i, elem := range sli {
		if cb(elem) {
			return i
		}
	}
	return -1
}

func IndexOfElemInSlice[E comparable](sli []E, v E) int {
	for i, elem := range sli {
		if elem == v {
			return i
		}
	}
	return -1
}

func IndexOfIntInSlice(sli []int, v int) int {
	for i, elem := range sli {
		if elem == v {
			return i
		}
	}
	return -1
}
func DoesIntSliceContain(sli []int, v int) bool {
	for _, elem := range sli {
		if elem == v {
			return true
		}
	}
	return false
}

func GetFloat64sAverageExcludeZeroValuesCheckIfAllZero(sli []float64) (float64, bool) {
	var count int
	var sum float64
	for _, elem := range sli {
		if 0.0 != elem {
			count++
			sum += elem
		}
	}
	if 0 == count {
		return 0.0, true
	}
	return sum / float64(count), false
}
func GetFloat64sAverageExcludeZeroValues(sli []float64) float64 {
	var count int
	var sum float64
	for _, elem := range sli {
		if 0.0 != elem {
			count++
			sum += elem
		}
	}
	//if 0 == count {
	//	return math.NaN()
	//}
	return sum / float64(count)
}

func CheckFloat64sAreAllZeroValues(sli []float64) bool {
	for _, elem := range sli {
		if 0.0 != elem {
			return false
		}
	}
	return true
}
func GetElemInFloatSliceIfExists(sli []float64, ind int) (float64, bool) {
	if ind < len(sli) {
		return sli[ind], true
	}
	return 0, false
}
func GetNthLastElemInFloatSlice(sli []float64, nthLast int) float64 {
	return sli[len(sli)-nthLast]
}

func GetSmallestFloatGeXInSortedSlice(sli []float64, limit float64) (float64, bool) {
	ind := sort.SearchFloat64s(sli, limit)
	if ind == len(sli) {
		return 0, false
	}
	return sli[ind], true
}
func GetGreatestFloatLtXInSortedSlice(sli []float64, limit float64) (float64, bool) {
	ind := sort.SearchFloat64s(sli, limit)
	if 0 == ind {
		return 0, false
	}
	return sli[ind-1], true
}
func GetMinInFloat64Slice(sl []float64) float64 {
	retval := sl[0]
	for _, v := range sl {
		retval = math.Min(v, retval)
	}
	return retval
}
func GetMaxInFloat64Slice(sl []float64) float64 {
	retval := sl[0]
	for _, v := range sl {
		retval = math.Max(v, retval)
	}
	return retval
}
func FindMinInFloat64sExcludeZeroValues(sl []float64) float64 {
	retval := math.MaxFloat64
	for _, v := range sl {
		if 0.0 != v && v < retval {
			retval = v
		}
	}
	return retval
}
func FindMaxInFloat64sExcludeZeroValues(sl []float64) float64 {
	retval := math.Inf(-1)
	for _, v := range sl {
		if 0.0 != v && v > retval {
			retval = v
		}
	}
	return retval
}
