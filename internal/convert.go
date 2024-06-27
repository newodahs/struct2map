package internal

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	STRUCT_MAP_PRIMARY_TAGNAME   = "struct2map"
	STRUCT_MAP_TAG_OMIT          = "omitempty"     // for nil-able values only; if nil, don't add to map
	STRUCT_MAP_TAG_IGNORE_PARENT = "ignoreparents" // don't use any of the parent names above this item; parents still honored for items contained within this item
)

func ConvertAnyToString(val any) string {
	if val == nil {
		return ""
	}

	valOf := reflect.ValueOf(val)
	for {
		if valOf.Kind() == reflect.Pointer {
			if reflect.Value(valOf).IsNil() {
				return ""
			}
			valOf = valOf.Elem()
			continue
		}
		break
	}

	switch valOf.Kind() {
	case reflect.String:
		return valOf.String()
	case reflect.Complex64:
		return strconv.FormatComplex(valOf.Complex(), 'g', -1, 64)
	case reflect.Complex128:
		return strconv.FormatComplex(valOf.Complex(), 'g', -1, 128)
	case reflect.Float32:
		return strconv.FormatFloat(valOf.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(valOf.Float(), 'g', -1, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(valOf.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(valOf.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(valOf.Bool())
	}

	return fmt.Sprintf("%v", val) // we don't support a proper conversion but let's return /something/ and hope for the best...
}
