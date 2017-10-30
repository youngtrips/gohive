// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type OptionValue struct {
	Bool     bool
	String   string
	Int      int64
	StrArray []string
	IntArray []int64
	Set      bool
}

func assign(value interface{}) (returnValue OptionValue, err *GetOptError) {
	valType := reflect.TypeOf(value).String()
	var e error

	// mmm ...there should be an easier way
	switch valType {
	case "string":
		returnValue.String = value.(string)
	case "bool":
		returnValue.Bool = value.(bool)
	case "int":
		returnValue.Int = int64(value.(int))
	case "int64":
		returnValue.Int = value.(int64)
	case "[]string":
		returnValue.StrArray = value.([]string)
	case "[]int":
		var ints []int = value.([]int)
		long_ints := make([]int64, len(ints))
		for i, integer := range ints {
			long_ints[i] = int64(integer)
		}
		returnValue.IntArray = long_ints
	case "[]int64":
		returnValue.IntArray = value.([]int64)
	default:
		e = errors.New("Couldn't assign value of type '" + valType + "'")
	}

	if e == nil {
		returnValue.Set = true
	} else {
		err = &GetOptError{OptionValueError, "Conversion Error: " + e.Error()}
	}

	return

}

func assignValue(referenceValue interface{}, value string) (returnValue OptionValue, err *GetOptError) {
	valType := reflect.TypeOf(referenceValue).String()
	var e error

	switch valType {
	case "string":
		returnValue.String = value
	case "bool":
		returnValue.Bool, e = strconv.ParseBool(value)
	case "int":
		fallthrough
	case "int64":
		returnValue.Int, e = strconv.ParseInt(value, 10, 64)
	case "[]string":
		returnValue.StrArray = strings.Split(value, ",")
	case "[]int":
		fallthrough
	case "[]int64":
		stringArray := strings.Split(value, ",")
		returnValue.IntArray = make([]int64, len(stringArray))
		for i, value := range stringArray {
			returnValue.IntArray[i], e = strconv.ParseInt(value, 10, 64)
		}
	default:
		e = errors.New("Couldn't convert '" + value + "' to type '" + valType + "'")
	}

	if e == nil {
		returnValue.Set = true
	} else {
		err = &GetOptError{OptionValueError, "Conversion Error: " + e.Error()}
	}

	return
}
