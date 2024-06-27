# struct2map: Convert golang strutures to a single, flattened map #

Primary usecases:
 1. Structures that you may wish to itterate over in an key-value type fashion
 2. Strucutres that you wish to pull data out of without having to write code paths for each value (just use the namespace-like keys to get your data)

## Usage ##
```
func ConvertStruct(obj any, opts ...StructConvertOpts) map[string]any
```
Takes a structure `obj` and returns a `map[string]any` where the map keys are the field names in the structure and the values are the field values or nil on error (ex: empty struct passed; not a struct passed).

The `opts` argument is a list of optional modifying options you may pass; currently the only supported modifiers are as follows:
```
	STRUCT_CONVERT_NOOP               // does nothing
	STRUCT_CONVERT_MAPKEY_TOLOWER     // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to lowercase for the map output
	STRUCT_CONVERT_MAPKEY_TOUPPER     // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to UPPERCASE for the map output
	STRUCT_CONVERT_MAPKEY_CAMELCASE   // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to CamelCase for the map output
	STRUCT_CONVERT_MAPKEY_LOWERCAMEL  // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to lowerCamelCase for the map output
	STRUCT_CONVERT_MAPKEY_SNAKE       // ignores the struct2map tag name (if set) and converts the STRUCT fieldname to snake_case for the map output
```
These modifiers are meant for a case where you cannot decorate a structure with the `struct2map` tag (see below), for example if the struct is from a third-party library you cannot or do not wish to modify.  All of the `STRUCT_CONVERT_MAPKEY_*` modifier options are BEST EFFORT only AND are MUTUALLY EXCLUSIVE to one another; last one passed should win, but please don't pass more than one...

## Notes ##

Most of the basic types at this point are supported for the map values, including nested/embedded structures, maps, slices, etc...

Pointers are always dereferenced when storing the values (avoid storing the pointer address).

Only operates on exported fields in the strucuture; non-exported fields are ignored.

Output maps are keyed by either the exported field name directly OR by the use of the `struct2map` tag to specify a name. If the name is specified as `-`, then the field is treated as not-exported.

In the case of slices, maps, or other embedded/nested structures, the output maps keys are "namespaced" in the following ways:
 * **Structures**: The embedded/nested structure name will prepend the inner fields as `[parentFieldName].[childFieldName] => value`.
 * **Maps**: Data pulled from maps will appears as `[mapFieldName].[mapKeyToString] => [value]`.
   * A nil pointer map key will ultimately bubble up as the string constant `DEFAULT_SUBKEY_STRING` wrapped in squre brackets (ex. with no convert options): `[emptyKey]`.
   * If the map key is otherwise unable to be directly converted to a string, we make a best effort via the `%v` format specifier with `fmt.Sprintf`.
   * In the event the map key is a `float` (of any type) or `complex64`/`complex128`, the conversion function to string uses the '`g`' modifier with a precision of -1 (see notes on https://pkg.go.dev/strconv#FormatFloat and https://pkg.go.dev/strconv#FormatComplex).
   * In all cases, these keys are also subject to the conversion options (above).
 * **Slices**: Data pulled form slices will appear as `[sliceFieldName].[sliceIndex] => [value]`.

As the amount of nesting increases, so does the namespacing; for example:
```
type someStruct struct {
    TopLevel bool
    InnerStruct anotherStruct
}

type anotherStruct struct {
    SomeMap map[string]int
}

...

testStruct := someStruct{
    TopLevel: true,
    InnerStruct: anotherStruct {
        SomeMap: map[string]int{"test": 1}
    }
}

outputMap := ConvertStruct(testStruct)
```
The key for an item on the `SomeMap` map in the `InnerStruct` would appear in the output map as the following key: `InnerStruct.SomeMap.test => 1`.

Additional tag options include (comma-separated, after the name):
 * `omitempty` - nil-able (and only nil-able) types are not added to the output map if set to nil.
 * `ignoreparents` - ignores all of the parents (prefixes) above the current position of nested fields, effectively flattening the keys (to a degree; beware of potential output map key conflicts when using this).

For `ignoreparents`, given the same `someStruct` example above, if the `SomeMap` field were to have `ignoreparents` then it would be keyed as the following in the output map: `SomeMap.test => [value]` (loss of the `InnerStruct` prefix).
