package utils

import (
	"fmt"
	"os"
	"reflect"
)

var Verbose bool

func DebugPrintln(message string) {
	if Verbose {
		StdErrorPrintln(message)
	}
}

func StdErrorPrintln(message string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", a...)
}

func PrintStructFields(s interface{}) {
	if s == nil {
		fmt.Println("nil value provided")

		return
	}

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			fmt.Println("nil pointer provided")

			return
		}

		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Println("provided value is not a struct")

		return
	}

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i)

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			fmt.Printf("%s: nil\n", fieldName)
		} else {
			fmt.Printf("%s: %v\n", fieldName, fieldValue.Interface())
		}
	}
}
