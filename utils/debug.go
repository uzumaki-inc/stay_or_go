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

func PrintStructFields(structObj interface{}) {
	if structObj == nil {
		fmt.Println("nil value provided")

		return
	}

	val := reflect.ValueOf(structObj)
	typ := reflect.TypeOf(structObj)

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

	for i := range make([]struct{}, val.NumField()) {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i)

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			fmt.Printf("%s: nil\n", fieldName)
		} else {
			fmt.Printf("%s: %v\n", fieldName, fieldValue.Interface())
		}
	}
}
