package golangrs

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInsertFloatToSortedSlice(t *testing.T) {
	buf := []float64{1, 2, 3, 4, 5}
	InsertFloatToSortedSlice(&buf, 2.5)
	fmt.Println(buf)
	if !reflect.DeepEqual(buf, []float64{1, 2, 2.5, 3, 4, 5}) {
		t.Errorf(`InsertFloatToSortedSlice failure`)
	}
	InsertFloatToSortedSlice(&buf, 4.5)
	fmt.Println(buf)
	if !reflect.DeepEqual(buf, []float64{1, 2, 2.5, 3, 4, 4.5, 5}) {
		t.Errorf(`InsertFloatToSortedSlice failure`)
	}
	InsertFloatToSortedSlice(&buf, 5.5)
	fmt.Println(buf)
	if !reflect.DeepEqual(buf, []float64{1, 2, 2.5, 3, 4, 4.5, 5, 5.5}) {
		t.Errorf(`InsertFloatToSortedSlice failure`)
	}
	InsertFloatToSortedSlice(&buf, 0.5)
	fmt.Println(buf)
	if !reflect.DeepEqual(buf, []float64{0.5, 1, 2, 2.5, 3, 4, 4.5, 5, 5.5}) {
		t.Errorf(`InsertFloatToSortedSlice failure`)
	}
}

func TestSwapFloat64sIfNotInSortedOrder(t *testing.T) {
	bar := 8.9
	foo := 6.7
	SwapFloat64sIfNotInSortedOrder(&bar, &foo)
	fmt.Println(bar, foo)
	if bar != 6.7 || foo != 8.9 {
		t.Errorf(`SwapFloat64sIfNotInSortedOrder failure`)
	}
	bar = 3.4
	foo = 3.5
	SwapFloat64sIfNotInSortedOrder(&bar, &foo)
	fmt.Println(bar, foo)
	if bar != 3.4 || foo != 3.5 {
		t.Errorf(`SwapFloat64sIfNotInSortedOrder failure`)
	}
}
