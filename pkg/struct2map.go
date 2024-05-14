package struct2map

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/newodahs/struct2map/internal"
)

type StructConvertOpts uint

// For MAPKEY opts - they are mutually exclusive and the last one wins (prior options will basically be ignored if passed together; don't do this...)
const (
	STRUCT_CONVERT_NOOP              StructConvertOpts = iota // does nothing
	STRUCT_CONVERT_MAPKEY_TOLOWER                             // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to lowercase for the map output
	STRUCT_CONVERT_MAPKEY_TOUPPER                             // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to UPPERCASE for the map output
	STRUCT_CONVERT_MAPKEY_CAMELCASE                           // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to CamelCase for the map output
	STRUCT_CONVERT_MAPKEY_LOWERCAMEL                          // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to lowerCamelCase for the map output
	STRUCT_CONVERT_MAPKEY_SNAKE                               // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to snake_case for the map output
)

// Takes a structure (obj) and turns it into a single, flat map; allows passing of various options (see StructConvertOpts constants)
//
// The structure field names come the map keys while the field values become the values for the map.
// The map key names may be altered by using the struct2map tag on the structure field
//
// Nesting of other structures, slices, or maps within top level (passed) structure result in keys
// that represent in a parent-child namespace like format of [parentField].[childField];
// see README documentation for further notes on this
//
// Additional structure tag options include omitempty to omit nil-able fields from the map
// and ignoreparents to ignore the prior parent namespace prefixes at that point
//
// Returns: map[string]any that is representative of the passed structure or nil on error (ex: empty struct passed; not a struct passed)
func ConvertStruct(obj any, opts ...StructConvertOpts) map[string]any {
	var nameMod func(string) string

	for _, opt := range opts {
		switch opt {
		case STRUCT_CONVERT_MAPKEY_TOLOWER:
			nameMod = strings.ToLower
		case STRUCT_CONVERT_MAPKEY_TOUPPER:
			nameMod = strings.ToUpper
		case STRUCT_CONVERT_MAPKEY_CAMELCASE:
			nameMod = strcase.ToCamel
		case STRUCT_CONVERT_MAPKEY_LOWERCAMEL:
			nameMod = strcase.ToLowerCamel
		case STRUCT_CONVERT_MAPKEY_SNAKE:
			nameMod = strcase.ToSnake
		}
	}

	return structToMap(nameMod, "", obj)
}

func structToMap(nameModFunc func(string) string, parentName string, obj any) map[string]any {
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
		omitempty := false
		ignoreParents := false
		actualFieldName := objType.Field(pos).Name
		mapKeyName, ok := objType.Field(pos).Tag.Lookup(internal.STRUCT_MAP_PRIMARY_TAGNAME)
		if !ok {
			mapKeyName = actualFieldName //no tag, just take the field name
		} else {
			//proc the tag information
			fieldSplit := strings.Split(mapKeyName, ",")
			mapKeyName = fieldSplit[0] //fieldname is always pos 0 for us...

			// field should not be exported; ignore everything else after that as it's moot
			if mapKeyName == "-" {
				continue STRUCT_MEMBER_PROC
			}

			for fIdx, fVal := range fieldSplit {
				if fIdx < 1 {
					continue
				}

				switch fVal {
				case internal.STRUCT_MAP_TAG_IGNORE_PARENT:
					ignoreParents = true
				case internal.STRUCT_MAP_TAG_OMIT:
					omitempty = true
				}
			}

			// before we go, reset our key name to the actual field name if modifier function was passed to us...
			// we do this here because we have to process other tags (ignoreparents, omitemtpy) even when a modifier
			// is passed...
			if nameModFunc != nil {
				mapKeyName = actualFieldName
			}
		}

		// if we have a parent name, prepend it here (if not ignored)
		if ignoreParents {
			parentName = ""
		}

		fieldToMap(ret, parentName, mapKeyName, objValue.Field(pos), omitempty, nameModFunc)
	}

	return ret
}

func fieldToMap(dest map[string]any, parentKeyName, mapKeyName string, workingField reflect.Value, omitEmpty bool, nameModFunc func(string) string) {
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

	// if we were passed a valid name modifying function, call it upfront
	keyName := mapKeyName
	if nameModFunc != nil {
		keyName = nameModFunc(mapKeyName)
	}

	// setup the actual keyname if there is a parent
	if parentKeyName != "" {
		keyName = fmt.Sprintf("%s.%s", parentKeyName, keyName)
	}

	if !workingField.IsValid() {
		if !omitEmpty {
			dest[keyName] = nil
		}
		return
	}

	switch workingField.Kind() {
	case reflect.Struct:
		// start the process on a new struct
		for k, v := range structToMap(nameModFunc, keyName, workingField.Interface()) {
			dest[k] = v
		}
	case reflect.Map:
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
				for k, v := range structToMap(nameModFunc, keyName, mapVal.Interface()) {
					dest[k] = v
				}
			} else {
				dest[fmt.Sprintf("%s.%s", keyName, internal.ConvertAnyToString(mapItr.Key().Interface()))] = mapVal.Interface()
			}
		}
	case reflect.Slice:
		if omitEmpty && workingField.IsNil() {
			return
		}

		for idx := 0; idx < workingField.Len(); idx++ {
			sliceValue := workingField.Index(idx)
			if sliceValue.Kind() == reflect.Pointer {
				sliceValue = sliceValue.Elem()
			}

			innerSliceName := fmt.Sprintf("%s.%d", keyName, idx)
			if sliceValue.Kind() == reflect.Struct {
				for k, v := range structToMap(nameModFunc, innerSliceName, sliceValue.Interface()) {
					dest[k] = v
				}
			} else {
				dest[innerSliceName] = sliceValue.Interface()
			}
		}
	default:
		if omitEmpty && workingField.Kind() == reflect.Interface && workingField.IsNil() {
			return
		}

		dest[keyName] = workingField.Interface()
	}
}
