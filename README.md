# struct2map: Convert golang strutures to a single, flattened map #

Primary usecases:
 1. Structures that you may wish to itterate over in an key-value type fashion
 2. Strucutres that you wish to pull data out of without having to write code paths for each value (just use the namespace-like keys to get your data)

## Usage and Notes ##
```
func StructToMap(obj any) map[string]any
```
Takes a structure `obj` and returns a `map[string]any` where the map keys are the field names in the structure and the values are the field values or nil on error (ex: empty struct passed; not a struct passed).

Most of the basic types at this point are supported, including nested/embedded structures, maps, slices, etc...

Pointers are always dereferenced when storing the values (avoid storing the pointer address).

Only operates on exported fields in the strucuture; non-exported fields are ignored.

Output maps are keyed by either the exported field name directly OR by the use of the `struct2map` tag to specify a name. If the name is specified as `-`, then the field is treated as not-exported.

In the case of slices, maps, or other embedded/nested structures, the output maps keys are "namespaced" in the following ways:
 * **Structures**: The embedded/nested structure name will prepend the inner fields as `[parentFieldName].[childFieldName] => value`.
 * **Maps**: Data pulled from maps will appears as `[mapFieldName].[mapKeyToString] => [value]`.
   * If the map key is unable to be directly converted to a string, the key will come back as the type name as seen via reflection.
 * **Slices**: Data pulled form slices will appear as `[sliceFieldName].[sliceIndex] => [value]`.

As the amount of nesting increases, so does the namespacing; in example:
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

outputMap := StructToMap("", testStruct)
```
The key for an item on the `SomeMap` map in the `InnerStruct` would appear in the output map as the following key: `InnerStruct.SomeMap.test => 1`.

Additional tag options include (comma-separated, after the name):
 * `omitempty` - nil-able (and only nil-able) types are not added to the output map if set to nil.
 * `ignoreparents` - ignores all of the parents (prefixes) above the current position of nested fields, effectively flattening the keys (to a degree; beware of potential output map key conflicts when using this).

For `ignoreparents`, given the same `someStruct` example above, if the `SomeMap` field were to have `ignoreparents` then it would be keyed as the following in the output map: `SomeMap.test => [value]` (loss of the `InnerStruct` prefix).