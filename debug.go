package golangrs

import (
	"fmt"
)

var mapForDebugPrintlnOnce map[string]bool = make(map[string]bool)

func DebugPrintlnOnce(str string) {
	_, ok := mapForDebugPrintlnOnce[str]
	if ok {
		return
	}
	fmt.Println(str)
	mapForDebugPrintlnOnce[str] = true
}
