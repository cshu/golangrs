package golangrs

import (
	"fmt"
	"reflect"
	"testing"
)

type Tup struct {
	I int
	S string
}

func TestRemoveFromSlice(t *testing.T) {
	fmt.Println(`TestRemoveFromSlice`)
	var s1 []int
	s1 = []int{1, 2, 4, 3}
	fmt.Println(s1, len(s1), cap(s1))
	RemoveFromSlice(&s1, 2)
	fmt.Println(s1, len(s1), cap(s1))
	var s2 []Tup
	s2 = []Tup{Tup{1, `111`}, Tup{2, `222`}, Tup{3, `333`}}
	RemoveFromSlice(&s2, 1)
	fmt.Println(s2, len(s2), cap(s2))
}

func TestSortFloat64sRetainUnique(t *testing.T) {
	fmt.Println(`TestSortFloat64sRetainUnique`)
	var s1 []float64
	s1 = []float64{1.1, 2.2, 4.4, 3.3}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 2.2, 1.1, 4.4}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1, 2.2, 4.4}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 2.2, 4.4, 4.4}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1, 2.2, 4.4}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 1.1, 2.2, 3.3, 4.4, 4.4}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 1.1, 2.2, 2.2, 2.2, 2.2, 3.3, 4.4, 4.4}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 1.1, 1.1, 1.1}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1, 1.1}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
	s1 = []float64{1.1}
	SortFloat64sRetainUnique(&s1)
	fmt.Println(s1, len(s1), cap(s1))
	if !reflect.DeepEqual(s1, []float64{1.1}) {
		t.Errorf(`SortFloat64sRetainUnique failure`)
	}
}
