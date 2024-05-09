package struct2map

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/newodahs/struct2map/internal"
)

// Takes a structure (obj) and returns a map[string]any where the string keys are the fields of the structure and the values are the field values
//
// returns nil on error (ex: empty struct passed; not a struct passed)
func StructToMap(obj any) map[string]any {
	return structToMap("", obj)
}

func structToMap(parentName string, obj any) map[string]any {
	if obj == nil {
		return nil
	}

	objValue := reflect.ValueOf(obj)

	for {
		if objValue.Kind() == reflect.Pointer {
			objValue = objValue.Elem()
			continue
		}
		break
	}

	if objValue.Kind() != reflect.Struct { // we only operate on structs
		return nil
	}

	ret := make(map[string]any)

	objType := objValue.Type()

	//rip over each structure member and process it into the map
STRUCT_MEMBER_PROC:
	for pos := 0; pos < objValue.NumField(); pos++ {
		if !objType.Field(pos).IsExported() {
			continue STRUCT_MEMBER_PROC
		}

		//proc the field tags (if any)
		mapKeyName, ok := objType.Field(pos).Tag.Lookup(internal.STRUCT_MAP_PRIMARY_TAGNAME)
		if !ok {
			mapKeyName = objType.Field(pos).Name
		}

		omitempty := false
		ignoreParents := false
		if fieldSplit := strings.Split(mapKeyName, ","); len(fieldSplit) > 1 {
			mapKeyName = fieldSplit[0] //fieldname is always pos 0

			for _, fVal := range fieldSplit[1:] {
				switch fVal {
				case internal.STRUCT_MAP_TAG_IGNORE_PARENT:
					ignoreParents = true
				case internal.STRUCT_MAP_TAG_OMIT:
					omitempty = true
				}
			}
		}

		// field should not be exported
		if mapKeyName == "-" {
			continue STRUCT_MEMBER_PROC
		}

		//if we have a parent name, prepend it here (if not ignored)
		if !ignoreParents && parentName != "" {
			mapKeyName = fmt.Sprintf("%s.%s", parentName, mapKeyName)
		}

		fieldToMap(ret, mapKeyName, objValue.Field(pos), omitempty)
	}

	return ret
}

func fieldToMap(dest map[string]any, mapKeyName string, workingField reflect.Value, omitEmpty bool) {
	for {
		if workingField.Kind() == reflect.Pointer {
			if omitEmpty && workingField.IsNil() {
				return
			}
			workingField = workingField.Elem()
			continue
		}
		break
	}

	switch workingField.Kind() {
	case reflect.Struct:
		// start the process on a new struct
		for k, v := range structToMap(mapKeyName, workingField.Interface()) {
			dest[k] = v
		}
	case reflect.Map:
		if !workingField.IsValid() {
			if !omitEmpty {
				dest[mapKeyName] = nil
			}
			return
		}

		if omitEmpty && workingField.IsNil() {
			return
		}

		mapItr := workingField.MapRange()
		for mapItr.Next() {
			mapVal := mapItr.Value()
			if mapVal.Kind() == reflect.Pointer {
				mapVal = mapVal.Elem()
			}

			if mapVal.Kind() == reflect.Struct {
				for k, v := range structToMap(mapKeyName, mapVal.Interface()) {
					dest[k] = v
				}
			} else {
				dest[fmt.Sprintf("%s.%s", mapKeyName, internal.ConvertAnyToString(mapItr.Key().Interface()))] = mapVal.Interface()
			}
		}
	case reflect.Slice:
		if !workingField.IsValid() {
			if !omitEmpty {
				dest[mapKeyName] = nil
			}
			return
		}

		if omitEmpty && workingField.IsNil() {
			return
		}

		for idx := 0; idx < workingField.Len(); idx++ {
			sliceValue := workingField.Index(idx)
			if sliceValue.Kind() == reflect.Pointer {
				sliceValue = sliceValue.Elem()
			}

			innerSliceName := fmt.Sprintf("%s.%d", mapKeyName, idx)
			if sliceValue.Kind() == reflect.Struct {
				for k, v := range structToMap(innerSliceName, sliceValue.Interface()) {
					dest[k] = v
				}
			} else {
				dest[innerSliceName] = sliceValue.Interface()
			}
		}
	default:
		if !workingField.IsValid() {
			if !omitEmpty {
				dest[mapKeyName] = nil
			}
			return
		}

		if omitEmpty && workingField.Kind() == reflect.Interface && workingField.IsNil() {
			return
		}

		dest[mapKeyName] = workingField.Interface()
	}
}
