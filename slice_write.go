package golangrs

import (
	"sort"
)

func InsertFloatToSlice(sli *[]float64, v float64, ind int) {
	lenSli := len(*sli)
	if ind == lenSli {
		*sli = append(*sli, v)
		return
	}
	*sli = append(*sli, (*sli)[lenSli-1])
	copy((*sli)[ind+1:lenSli], (*sli)[ind:lenSli-1])
	(*sli)[ind] = v
}
func InsertFloatToSortedSlice(sli *[]float64, v float64) {
	ind := sort.SearchFloat64s(*sli, v)
	InsertFloatToSlice(sli, v, ind)
}
func RemoveFloatInSortedSlice(sli *[]float64, v float64) {
	ind := sort.SearchFloat64s(*sli, v)
	*sli = append((*sli)[:ind], (*sli)[ind+1:]...)
}
func ReplaceFloatInSortedSliceAndSortAgain(sli []float64, old float64, nouveau float64) {
	ind := sort.SearchFloat64s(sli, old)
	sli[ind] = nouveau
	sort.Float64s(sli)
}

// note this func assumes len is at least 1
func SortFloat64sRetainUnique(sli *[]float64) {
	sort.Float64s(*sli)
	lastCopiedIdx := 0
	nextComparIdx := 1
	for ; nextComparIdx != len(*sli); nextComparIdx++ {
		if (*sli)[lastCopiedIdx] != (*sli)[nextComparIdx] {
			lastCopiedIdx++
			(*sli)[lastCopiedIdx] = (*sli)[nextComparIdx]
		}
	}
	*sli = (*sli)[:lastCopiedIdx+1]
}

func RemoveFromSlice[S ~[]E, E any](sli *S, idx int) {
	*sli = append((*sli)[:idx], (*sli)[idx+1:]...)
}

func InsertToSlice[S ~[]E, E any](sli *S, v E, ind int) {
	lenSli := len(*sli)
	if ind == lenSli {
		*sli = append(*sli, v)
		return
	}
	*sli = append(*sli, (*sli)[lenSli-1])
	copy((*sli)[ind+1:lenSli], (*sli)[ind:lenSli-1])
	(*sli)[ind] = v
}
