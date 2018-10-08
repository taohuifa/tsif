package main

import (
	"fmt"
	"reflect"

	"github.com/tsif/utils"
)

func main() {
	fmt.Println("hello world")

	kinds := []reflect.Kind{
		reflect.Bool,
		reflect.Float32,
		reflect.Float64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.String,
	}

	for i := 0; i < len(kinds); i++ {
		str := "1234567"
		kind := kinds[i]
		v := utils.ChangeType(str, kind)
		fmt.Println(str, kind, v, reflect.TypeOf(v))
	}

}
