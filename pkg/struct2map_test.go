package struct2map

import (
	"log"
	"testing"
)

type simpleTestStruct struct {
	notExportedField           int `struct2map:"notExported"`
	NegatedExportField         int `struct2map:"-"`
	RegularFieldNoTag          int
	RegularFieldNameTag        int   `struct2map:"regularField"`
	RegularFieldOmitEmpty      *int  `struct2map:"regularFieldOmitEmpty,omitempty"`
	RegularFieldPointerPointer **int `struct2map:"regularFieldPointerPointer"`
}

type complexTestStruct struct {
	TopLevelField              bool                         `struct2map:"topLevelBool"`
	SliceField                 []int                        `struct2map:"sliceField"`
	SliceFieldPtrVal           []*int                       `struct2map:"sliceFieldPtrVal"`
	MapFieldStrKey             map[string]int               `struct2map:"mapFieldStrKey"`
	MapFieldStrKeyPtrVal       map[string]*int              `struct2map:"mapFieldStrKeyPtrVal"`
	MapFieldIntKey             map[int]string               `struct2map:"mapFieldIntKey"`
	MapFieldStrKeyStructVal    map[string]simpleTestStruct  `struct2map:"mapFieldStrKeyStructVal"`
	MapFieldStrKeyStructPtrVal map[string]*simpleTestStruct `struct2map:"mapFieldStrKeyStructPtrVal"`
	MapFieldPointerKey         map[*string]string           `struct2map:"mapFieldPointerKey"`
}

type complexTestStructEmbed struct {
	ComplexTestStruct complexTestStruct
	AnonStruct        struct {
		RegStruct    simpleTestStruct
		RegStructPtr *simpleTestStruct
	} `struct2map:"anonStruct"`
}

type embeddedStruct struct {
	TopLevelValue              bool             `struct2map:"topLevelBool"`
	notExportedStruct          simpleTestStruct `struct2map:"notExportedStruct"`
	NegatedExportedStruct      simpleTestStruct `struct2map:"-"`
	RegularExportStructNoTag   simpleTestStruct
	RegularFieldNameTag        simpleTestStruct   `struct2map:"regularStruct"`
	RegularFieldOmitEmpty      *simpleTestStruct  `struct2map:"regularStructOmitEmpty,omitempty"`
	RegularFieldPointerPointer **simpleTestStruct `struct2map:"regularStructPointerPointer"`
	AnonStruct                 struct {
		RegStruct    simpleTestStruct
		RegStructPtr *simpleTestStruct
	} `struct2map:"anonStruct"`
	AnonStructPtr *struct {
		RegStruct    simpleTestStruct
		RegStructPtr *simpleTestStruct
	} `struct2map:"anonStructPtr"`
}

type flattenStruct struct {
	TopLevelValue       bool `struct2map:"topLevelBool"`
	AnonContainedStruct struct {
		ContainedStruct *simpleTestStruct `struct2map:"structIgnoreParent,ignoreparents"`
	} `struct2map:"anonContained"`
}

// test case set for our most basic cases, more complicated and/or corner cases should go elsewhere
func Test_RegularCases(t *testing.T) {
	simpleInt := 1
	simpleIntPtr := &simpleInt

	testStruct := simpleTestStruct{
		notExportedField:           simpleInt,
		NegatedExportField:         simpleInt,
		RegularFieldNoTag:          simpleInt,
		RegularFieldNameTag:        simpleInt,
		RegularFieldOmitEmpty:      &simpleInt,
		RegularFieldPointerPointer: &simpleIntPtr,
	}
	testStructPtr := &testStruct

	testStrKey := "testKey"

	testSet := []struct {
		Name          string
		TestStructure any
		ExpectedMap   map[string]any
		SkipTest      bool
	}{
		{
			Name:          "simpleTestStruct validation",
			TestStructure: testStruct,
			ExpectedMap:   map[string]any{"RegularFieldNoTag": 1, "regularField": 1, "regularFieldOmitEmpty": simpleInt, "regularFieldPointerPointer": simpleInt},
		},
		{
			Name:          "pointer of simpleTestStruct validation",
			TestStructure: testStructPtr,
			ExpectedMap:   map[string]any{"RegularFieldNoTag": 1, "regularField": 1, "regularFieldOmitEmpty": simpleInt, "regularFieldPointerPointer": simpleInt},
		},
		{
			Name: "simpleTestStruct validation RegularFieldOmitEmpty nil to omit RegularFieldPointerPointer nil and not omitted",
			TestStructure: simpleTestStruct{
				notExportedField:    simpleInt,
				NegatedExportField:  simpleInt,
				RegularFieldNoTag:   simpleInt,
				RegularFieldNameTag: simpleInt,
			},
			ExpectedMap: map[string]any{"RegularFieldNoTag": simpleInt, "regularField": simpleInt, "regularFieldPointerPointer": nil},
		},
		{
			Name: "complexTestStruct validation",
			TestStructure: complexTestStruct{
				TopLevelField:              true,
				SliceField:                 []int{simpleInt, simpleInt, simpleInt},
				SliceFieldPtrVal:           []*int{&simpleInt, &simpleInt, &simpleInt},
				MapFieldStrKey:             map[string]int{"field-a": simpleInt, "field-b": simpleInt},
				MapFieldStrKeyPtrVal:       map[string]*int{"field-a-ptr": &simpleInt, "field-b-ptr": &simpleInt},
				MapFieldIntKey:             map[int]string{1: "test1", 2: "test2"},
				MapFieldStrKeyStructVal:    map[string]simpleTestStruct{"simpleStruct1": *testStructPtr},
				MapFieldStrKeyStructPtrVal: map[string]*simpleTestStruct{"simpleStructPtr1": testStructPtr},
				MapFieldPointerKey:         map[*string]string{nil: "testing1", &testStrKey: "testing2"},
			},
			ExpectedMap: map[string]any{
				"mapFieldIntKey.1":                                      "test1",
				"mapFieldIntKey.2":                                      "test2",
				"mapFieldStrKey.field-a":                                1,
				"mapFieldStrKey.field-b":                                1,
				"mapFieldStrKeyPtrVal.field-a-ptr":                      1,
				"mapFieldStrKeyPtrVal.field-b-ptr":                      1,
				"mapFieldStrKeyStructPtrVal.RegularFieldNoTag":          1,
				"mapFieldStrKeyStructPtrVal.regularField":               1,
				"mapFieldStrKeyStructPtrVal.regularFieldOmitEmpty":      1,
				"mapFieldStrKeyStructPtrVal.regularFieldPointerPointer": 1,
				"mapFieldStrKeyStructVal.RegularFieldNoTag":             1,
				"mapFieldStrKeyStructVal.regularField":                  1,
				"mapFieldStrKeyStructVal.regularFieldOmitEmpty":         1,
				"mapFieldStrKeyStructVal.regularFieldPointerPointer":    1,
				"sliceField.0":                  1,
				"sliceField.1":                  1,
				"sliceField.2":                  1,
				"sliceFieldPtrVal.0":            1,
				"sliceFieldPtrVal.1":            1,
				"sliceFieldPtrVal.2":            1,
				"topLevelBool":                  true,
				"mapFieldPointerKey.[emptyKey]": "testing1",
				"mapFieldPointerKey.testKey":    "testing2",
			},
		},
		{
			Name: "embeddedStruct Pointer validation",
			TestStructure: embeddedStruct{
				TopLevelValue: true,
				notExportedStruct: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				NegatedExportedStruct: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularExportStructNoTag: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularFieldNameTag: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularFieldOmitEmpty:      testStructPtr,
				RegularFieldPointerPointer: &testStructPtr,
				AnonStruct: struct {
					RegStruct    simpleTestStruct
					RegStructPtr *simpleTestStruct
				}{
					RegStruct:    *testStructPtr,
					RegStructPtr: testStructPtr,
				},
				AnonStructPtr: &struct {
					RegStruct    simpleTestStruct
					RegStructPtr *simpleTestStruct
				}{
					RegStruct:    testStruct,
					RegStructPtr: testStructPtr,
				},
			},
			ExpectedMap: map[string]any{
				"RegularExportStructNoTag.RegularFieldNoTag":             1,
				"RegularExportStructNoTag.regularField":                  1,
				"RegularExportStructNoTag.regularFieldOmitEmpty":         1,
				"RegularExportStructNoTag.regularFieldPointerPointer":    1,
				"anonStruct.RegStruct.RegularFieldNoTag":                 1,
				"anonStruct.RegStruct.regularField":                      1,
				"anonStruct.RegStruct.regularFieldOmitEmpty":             1,
				"anonStruct.RegStruct.regularFieldPointerPointer":        1,
				"anonStruct.RegStructPtr.RegularFieldNoTag":              1,
				"anonStruct.RegStructPtr.regularField":                   1,
				"anonStruct.RegStructPtr.regularFieldOmitEmpty":          1,
				"anonStruct.RegStructPtr.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStruct.RegularFieldNoTag":              1,
				"anonStructPtr.RegStruct.regularField":                   1,
				"anonStructPtr.RegStruct.regularFieldOmitEmpty":          1,
				"anonStructPtr.RegStruct.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStructPtr.RegularFieldNoTag":           1,
				"anonStructPtr.RegStructPtr.regularField":                1,
				"anonStructPtr.RegStructPtr.regularFieldOmitEmpty":       1,
				"anonStructPtr.RegStructPtr.regularFieldPointerPointer":  1,
				"regularStruct.RegularFieldNoTag":                        1,
				"regularStruct.regularField":                             1,
				"regularStruct.regularFieldOmitEmpty":                    1,
				"regularStruct.regularFieldPointerPointer":               1,
				"regularStructOmitEmpty.RegularFieldNoTag":               1,
				"regularStructOmitEmpty.regularField":                    1,
				"regularStructOmitEmpty.regularFieldOmitEmpty":           1,
				"regularStructOmitEmpty.regularFieldPointerPointer":      1,
				"regularStructPointerPointer.RegularFieldNoTag":          1,
				"regularStructPointerPointer.regularField":               1,
				"regularStructPointerPointer.regularFieldOmitEmpty":      1,
				"regularStructPointerPointer.regularFieldPointerPointer": 1,
				"topLevelBool": true,
			},
		},
		{
			Name: "pointer of embeddedStruct validation",
			TestStructure: &embeddedStruct{
				TopLevelValue: true,
				notExportedStruct: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				NegatedExportedStruct: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularExportStructNoTag: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularFieldNameTag: simpleTestStruct{
					notExportedField:           simpleInt,
					NegatedExportField:         simpleInt,
					RegularFieldNoTag:          simpleInt,
					RegularFieldNameTag:        simpleInt,
					RegularFieldOmitEmpty:      &simpleInt,
					RegularFieldPointerPointer: &simpleIntPtr,
				},
				RegularFieldOmitEmpty:      testStructPtr,
				RegularFieldPointerPointer: &testStructPtr,
				AnonStruct: struct {
					RegStruct    simpleTestStruct
					RegStructPtr *simpleTestStruct
				}{
					RegStruct:    testStruct,
					RegStructPtr: testStructPtr,
				},
				AnonStructPtr: &struct {
					RegStruct    simpleTestStruct
					RegStructPtr *simpleTestStruct
				}{
					RegStruct:    testStruct,
					RegStructPtr: testStructPtr,
				},
			},
			ExpectedMap: map[string]any{
				"RegularExportStructNoTag.RegularFieldNoTag":             1,
				"RegularExportStructNoTag.regularField":                  1,
				"RegularExportStructNoTag.regularFieldOmitEmpty":         1,
				"RegularExportStructNoTag.regularFieldPointerPointer":    1,
				"anonStruct.RegStruct.RegularFieldNoTag":                 1,
				"anonStruct.RegStruct.regularField":                      1,
				"anonStruct.RegStruct.regularFieldOmitEmpty":             1,
				"anonStruct.RegStruct.regularFieldPointerPointer":        1,
				"anonStruct.RegStructPtr.RegularFieldNoTag":              1,
				"anonStruct.RegStructPtr.regularField":                   1,
				"anonStruct.RegStructPtr.regularFieldOmitEmpty":          1,
				"anonStruct.RegStructPtr.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStruct.RegularFieldNoTag":              1,
				"anonStructPtr.RegStruct.regularField":                   1,
				"anonStructPtr.RegStruct.regularFieldOmitEmpty":          1,
				"anonStructPtr.RegStruct.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStructPtr.RegularFieldNoTag":           1,
				"anonStructPtr.RegStructPtr.regularField":                1,
				"anonStructPtr.RegStructPtr.regularFieldOmitEmpty":       1,
				"anonStructPtr.RegStructPtr.regularFieldPointerPointer":  1,
				"regularStruct.RegularFieldNoTag":                        1,
				"regularStruct.regularField":                             1,
				"regularStruct.regularFieldOmitEmpty":                    1,
				"regularStruct.regularFieldPointerPointer":               1,
				"regularStructOmitEmpty.RegularFieldNoTag":               1,
				"regularStructOmitEmpty.regularField":                    1,
				"regularStructOmitEmpty.regularFieldOmitEmpty":           1,
				"regularStructOmitEmpty.regularFieldPointerPointer":      1,
				"regularStructPointerPointer.RegularFieldNoTag":          1,
				"regularStructPointerPointer.regularField":               1,
				"regularStructPointerPointer.regularFieldOmitEmpty":      1,
				"regularStructPointerPointer.regularFieldPointerPointer": 1,
				"topLevelBool": true,
			},
		},
		{
			Name: "flattenStruct validation anonContained should not appear as a parent",
			TestStructure: flattenStruct{
				TopLevelValue: true,
				AnonContainedStruct: struct {
					ContainedStruct *simpleTestStruct "struct2map:\"structIgnoreParent,ignoreparents\""
				}{
					ContainedStruct: testStructPtr,
				},
			},
			ExpectedMap: map[string]any{
				"structIgnoreParent.RegularFieldNoTag":          1,
				"structIgnoreParent.regularField":               1,
				"structIgnoreParent.regularFieldOmitEmpty":      1,
				"structIgnoreParent.regularFieldPointerPointer": 1,
				"topLevelBool": true,
			},
		},
	}

	for _, curTest := range testSet {
		t.Run(curTest.Name, func(t *testing.T) {
			if curTest.SkipTest {
				t.Skipf("skipped '%s' due to SkipTest being set", curTest.Name)
			}

			genMap := ConvertStruct(curTest.TestStructure)

			//compare generated to expected
			tErr := false
			for k, v := range genMap {
				expVal, ok := curTest.ExpectedMap[k]
				if !ok {
					t.Errorf("failed to find '%s' from the generated map in the expected map", k)
					tErr = true
				}

				if v != expVal {
					t.Errorf("value stored for '%s' (%+v) in the generated map not the same as what is in the expected map (%+v)", k, v, expVal)
					tErr = true
				}
			}

			//compare expected to generated
			for k, v := range curTest.ExpectedMap {
				genVal, ok := genMap[k]
				if !ok {
					t.Errorf("failed to find '%s' from the expected map in the generated map", k)
					tErr = true
				}

				if v != genVal {
					t.Errorf("value stored for '%s' (%+v) in the expedcted map not the same as what is in the generated map (%+v)", k, v, genVal)
					tErr = true
				}
			}

			if tErr {
				log.Printf("Have: %+v", genMap)
				log.Printf("Want: %+v", curTest.ExpectedMap)
			}
		})
	}
}

// test case set for struct field name -> key name conversion options (ignoring the tag name)
func Test_MapKeyOptions(t *testing.T) {
	simpleInt := 1
	simpleIntPtr := &simpleInt
	simpleStruct := simpleTestStruct{
		notExportedField:           simpleInt,
		NegatedExportField:         simpleInt,
		RegularFieldNoTag:          simpleInt,
		RegularFieldNameTag:        simpleInt,
		RegularFieldOmitEmpty:      &simpleInt,
		RegularFieldPointerPointer: &simpleIntPtr,
	}

	testStrKey := "testKey"

	testStructPtr := &complexTestStructEmbed{
		ComplexTestStruct: complexTestStruct{
			TopLevelField:              true,
			SliceField:                 []int{1, 2, 3},
			SliceFieldPtrVal:           []*int{&simpleInt},
			MapFieldStrKey:             map[string]int{"test1": 1, "test2": 2, "test3": 3},
			MapFieldStrKeyPtrVal:       map[string]*int{"test1": &simpleInt},
			MapFieldIntKey:             map[int]string{1: "test1", 2: "test2"},
			MapFieldStrKeyStructVal:    map[string]simpleTestStruct{"test1": simpleStruct},
			MapFieldStrKeyStructPtrVal: map[string]*simpleTestStruct{"testPtr1": &simpleStruct},
			MapFieldPointerKey:         map[*string]string{nil: "testing1", &testStrKey: "testing2"},
		},
		AnonStruct: struct {
			RegStruct    simpleTestStruct
			RegStructPtr *simpleTestStruct
		}{
			RegStruct:    simpleStruct,
			RegStructPtr: &simpleStruct,
		},
	}

	testSet := []struct {
		Name          string
		TestStructure any
		ExpectedMap   map[string]any
		ConvertOpts   []StructConvertOpts
		SkipTest      bool
	}{
		{
			Name:          "to-lower mapkey validation",
			TestStructure: testStructPtr,
			ConvertOpts:   []StructConvertOpts{STRUCT_CONVERT_MAPKEY_TOLOWER},
			ExpectedMap: map[string]any{
				"anonstruct.regstruct.regularfieldnametag":                                1,
				"anonstruct.regstruct.regularfieldnotag":                                  1,
				"anonstruct.regstruct.regularfieldomitempty":                              1,
				"anonstruct.regstruct.regularfieldpointerpointer":                         1,
				"anonstruct.regstructptr.regularfieldnametag":                             1,
				"anonstruct.regstructptr.regularfieldnotag":                               1,
				"anonstruct.regstructptr.regularfieldomitempty":                           1,
				"anonstruct.regstructptr.regularfieldpointerpointer":                      1,
				"complexteststruct.mapfieldintkey.1":                                      "test1",
				"complexteststruct.mapfieldintkey.2":                                      "test2",
				"complexteststruct.mapfieldstrkey.test1":                                  1,
				"complexteststruct.mapfieldstrkey.test2":                                  2,
				"complexteststruct.mapfieldstrkey.test3":                                  3,
				"complexteststruct.mapfieldstrkeyptrval.test1":                            1,
				"complexteststruct.mapfieldstrkeystructptrval.regularfieldnametag":        1,
				"complexteststruct.mapfieldstrkeystructptrval.regularfieldnotag":          1,
				"complexteststruct.mapfieldstrkeystructptrval.regularfieldomitempty":      1,
				"complexteststruct.mapfieldstrkeystructptrval.regularfieldpointerpointer": 1,
				"complexteststruct.mapfieldstrkeystructval.regularfieldnametag":           1,
				"complexteststruct.mapfieldstrkeystructval.regularfieldnotag":             1,
				"complexteststruct.mapfieldstrkeystructval.regularfieldomitempty":         1,
				"complexteststruct.mapfieldstrkeystructval.regularfieldpointerpointer":    1,
				"complexteststruct.slicefield.0":                                          1,
				"complexteststruct.slicefield.1":                                          2,
				"complexteststruct.slicefield.2":                                          3,
				"complexteststruct.slicefieldptrval.0":                                    1,
				"complexteststruct.toplevelfield":                                         true,
				"complexteststruct.mapfieldpointerkey.[emptykey]":                         "testing1",
				"complexteststruct.mapfieldpointerkey.testkey":                            "testing2",
			},
		},
		{
			Name:          "to-upper mapkey validation",
			TestStructure: testStructPtr,
			ConvertOpts:   []StructConvertOpts{STRUCT_CONVERT_MAPKEY_TOUPPER},
			ExpectedMap: map[string]any{
				"ANONSTRUCT.REGSTRUCT.REGULARFIELDNAMETAG":                                1,
				"ANONSTRUCT.REGSTRUCT.REGULARFIELDNOTAG":                                  1,
				"ANONSTRUCT.REGSTRUCT.REGULARFIELDOMITEMPTY":                              1,
				"ANONSTRUCT.REGSTRUCT.REGULARFIELDPOINTERPOINTER":                         1,
				"ANONSTRUCT.REGSTRUCTPTR.REGULARFIELDNAMETAG":                             1,
				"ANONSTRUCT.REGSTRUCTPTR.REGULARFIELDNOTAG":                               1,
				"ANONSTRUCT.REGSTRUCTPTR.REGULARFIELDOMITEMPTY":                           1,
				"ANONSTRUCT.REGSTRUCTPTR.REGULARFIELDPOINTERPOINTER":                      1,
				"COMPLEXTESTSTRUCT.MAPFIELDINTKEY.1":                                      "test1",
				"COMPLEXTESTSTRUCT.MAPFIELDINTKEY.2":                                      "test2",
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEY.TEST1":                                  1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEY.TEST2":                                  2,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEY.TEST3":                                  3,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYPTRVAL.TEST1":                            1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTPTRVAL.REGULARFIELDNAMETAG":        1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTPTRVAL.REGULARFIELDNOTAG":          1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTPTRVAL.REGULARFIELDOMITEMPTY":      1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTPTRVAL.REGULARFIELDPOINTERPOINTER": 1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTVAL.REGULARFIELDNAMETAG":           1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTVAL.REGULARFIELDNOTAG":             1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTVAL.REGULARFIELDOMITEMPTY":         1,
				"COMPLEXTESTSTRUCT.MAPFIELDSTRKEYSTRUCTVAL.REGULARFIELDPOINTERPOINTER":    1,
				"COMPLEXTESTSTRUCT.SLICEFIELD.0":                                          1,
				"COMPLEXTESTSTRUCT.SLICEFIELD.1":                                          2,
				"COMPLEXTESTSTRUCT.SLICEFIELD.2":                                          3,
				"COMPLEXTESTSTRUCT.SLICEFIELDPTRVAL.0":                                    1,
				"COMPLEXTESTSTRUCT.TOPLEVELFIELD":                                         true,
				"COMPLEXTESTSTRUCT.MAPFIELDPOINTERKEY.[EMPTYKEY]":                         "testing1",
				"COMPLEXTESTSTRUCT.MAPFIELDPOINTERKEY.TESTKEY":                            "testing2",
			},
		},
		{
			Name:          "camelcase mapkey validation",
			TestStructure: testStructPtr,
			ConvertOpts:   []StructConvertOpts{STRUCT_CONVERT_MAPKEY_CAMELCASE},
			ExpectedMap: map[string]any{
				"AnonStruct.RegStruct.RegularFieldNameTag":                                1,
				"AnonStruct.RegStruct.RegularFieldNoTag":                                  1,
				"AnonStruct.RegStruct.RegularFieldOmitEmpty":                              1,
				"AnonStruct.RegStruct.RegularFieldPointerPointer":                         1,
				"AnonStruct.RegStructPtr.RegularFieldNameTag":                             1,
				"AnonStruct.RegStructPtr.RegularFieldNoTag":                               1,
				"AnonStruct.RegStructPtr.RegularFieldOmitEmpty":                           1,
				"AnonStruct.RegStructPtr.RegularFieldPointerPointer":                      1,
				"ComplexTestStruct.MapFieldIntKey.1":                                      "test1",
				"ComplexTestStruct.MapFieldIntKey.2":                                      "test2",
				"ComplexTestStruct.MapFieldStrKey.Test1":                                  1,
				"ComplexTestStruct.MapFieldStrKey.Test2":                                  2,
				"ComplexTestStruct.MapFieldStrKey.Test3":                                  3,
				"ComplexTestStruct.MapFieldStrKeyPtrVal.Test1":                            1,
				"ComplexTestStruct.MapFieldStrKeyStructPtrVal.RegularFieldNameTag":        1,
				"ComplexTestStruct.MapFieldStrKeyStructPtrVal.RegularFieldNoTag":          1,
				"ComplexTestStruct.MapFieldStrKeyStructPtrVal.RegularFieldOmitEmpty":      1,
				"ComplexTestStruct.MapFieldStrKeyStructPtrVal.RegularFieldPointerPointer": 1,
				"ComplexTestStruct.MapFieldStrKeyStructVal.RegularFieldNameTag":           1,
				"ComplexTestStruct.MapFieldStrKeyStructVal.RegularFieldNoTag":             1,
				"ComplexTestStruct.MapFieldStrKeyStructVal.RegularFieldOmitEmpty":         1,
				"ComplexTestStruct.MapFieldStrKeyStructVal.RegularFieldPointerPointer":    1,
				"ComplexTestStruct.SliceField.0":                                          1,
				"ComplexTestStruct.SliceField.1":                                          2,
				"ComplexTestStruct.SliceField.2":                                          3,
				"ComplexTestStruct.SliceFieldPtrVal.0":                                    1,
				"ComplexTestStruct.TopLevelField":                                         true,
				"ComplexTestStruct.MapFieldPointerKey.[EmptyKey]":                         "testing1",
				"ComplexTestStruct.MapFieldPointerKey.TestKey":                            "testing2",
			},
		},
		{
			Name:          "lower-camelcase mapkey validation",
			TestStructure: testStructPtr,
			ConvertOpts:   []StructConvertOpts{STRUCT_CONVERT_MAPKEY_LOWERCAMEL},
			ExpectedMap: map[string]any{
				"anonStruct.regStruct.regularFieldNameTag":                                1,
				"anonStruct.regStruct.regularFieldNoTag":                                  1,
				"anonStruct.regStruct.regularFieldOmitEmpty":                              1,
				"anonStruct.regStruct.regularFieldPointerPointer":                         1,
				"anonStruct.regStructPtr.regularFieldNameTag":                             1,
				"anonStruct.regStructPtr.regularFieldNoTag":                               1,
				"anonStruct.regStructPtr.regularFieldOmitEmpty":                           1,
				"anonStruct.regStructPtr.regularFieldPointerPointer":                      1,
				"complexTestStruct.mapFieldIntKey.1":                                      "test1",
				"complexTestStruct.mapFieldIntKey.2":                                      "test2",
				"complexTestStruct.mapFieldStrKey.test1":                                  1,
				"complexTestStruct.mapFieldStrKey.test2":                                  2,
				"complexTestStruct.mapFieldStrKey.test3":                                  3,
				"complexTestStruct.mapFieldStrKeyPtrVal.test1":                            1,
				"complexTestStruct.mapFieldStrKeyStructPtrVal.regularFieldNameTag":        1,
				"complexTestStruct.mapFieldStrKeyStructPtrVal.regularFieldNoTag":          1,
				"complexTestStruct.mapFieldStrKeyStructPtrVal.regularFieldOmitEmpty":      1,
				"complexTestStruct.mapFieldStrKeyStructPtrVal.regularFieldPointerPointer": 1,
				"complexTestStruct.mapFieldStrKeyStructVal.regularFieldNameTag":           1,
				"complexTestStruct.mapFieldStrKeyStructVal.regularFieldNoTag":             1,
				"complexTestStruct.mapFieldStrKeyStructVal.regularFieldOmitEmpty":         1,
				"complexTestStruct.mapFieldStrKeyStructVal.regularFieldPointerPointer":    1,
				"complexTestStruct.sliceField.0":                                          1,
				"complexTestStruct.sliceField.1":                                          2,
				"complexTestStruct.sliceField.2":                                          3,
				"complexTestStruct.sliceFieldPtrVal.0":                                    1,
				"complexTestStruct.topLevelField":                                         true,
				"complexTestStruct.mapFieldPointerKey.[emptyKey]":                         "testing1",
				"complexTestStruct.mapFieldPointerKey.testKey":                            "testing2",
			},
		},
		{
			Name:          "snakecase mapkey validation",
			TestStructure: testStructPtr,
			ConvertOpts:   []StructConvertOpts{STRUCT_CONVERT_MAPKEY_SNAKE},
			ExpectedMap: map[string]any{
				"anon_struct.reg_struct.regular_field_name_tag":                                      1,
				"anon_struct.reg_struct.regular_field_no_tag":                                        1,
				"anon_struct.reg_struct.regular_field_omit_empty":                                    1,
				"anon_struct.reg_struct.regular_field_pointer_pointer":                               1,
				"anon_struct.reg_struct_ptr.regular_field_name_tag":                                  1,
				"anon_struct.reg_struct_ptr.regular_field_no_tag":                                    1,
				"anon_struct.reg_struct_ptr.regular_field_omit_empty":                                1,
				"anon_struct.reg_struct_ptr.regular_field_pointer_pointer":                           1,
				"complex_test_struct.map_field_int_key.1":                                            "test1",
				"complex_test_struct.map_field_int_key.2":                                            "test2",
				"complex_test_struct.map_field_str_key.test_1":                                       1,
				"complex_test_struct.map_field_str_key.test_2":                                       2,
				"complex_test_struct.map_field_str_key.test_3":                                       3,
				"complex_test_struct.map_field_str_key_ptr_val.test_1":                               1,
				"complex_test_struct.map_field_str_key_struct_ptr_val.regular_field_name_tag":        1,
				"complex_test_struct.map_field_str_key_struct_ptr_val.regular_field_no_tag":          1,
				"complex_test_struct.map_field_str_key_struct_ptr_val.regular_field_omit_empty":      1,
				"complex_test_struct.map_field_str_key_struct_ptr_val.regular_field_pointer_pointer": 1,
				"complex_test_struct.map_field_str_key_struct_val.regular_field_name_tag":            1,
				"complex_test_struct.map_field_str_key_struct_val.regular_field_no_tag":              1,
				"complex_test_struct.map_field_str_key_struct_val.regular_field_omit_empty":          1,
				"complex_test_struct.map_field_str_key_struct_val.regular_field_pointer_pointer":     1,
				"complex_test_struct.slice_field.0":                                                  1,
				"complex_test_struct.slice_field.1":                                                  2,
				"complex_test_struct.slice_field.2":                                                  3,
				"complex_test_struct.slice_field_ptr_val.0":                                          1,
				"complex_test_struct.top_level_field":                                                true,
				"complex_test_struct.map_field_pointer_key.[empty_key]":                              "testing1",
				"complex_test_struct.map_field_pointer_key.test_key":                                 "testing2",
			},
		},
	}

	for _, curTest := range testSet {
		t.Run(curTest.Name, func(t *testing.T) {
			if curTest.SkipTest {
				t.Skipf("skipped '%s' due to SkipTest being set", curTest.Name)
			}

			genMap := ConvertStruct(curTest.TestStructure, curTest.ConvertOpts...)

			//compare generated to expected
			for k, v := range genMap {
				expVal, ok := curTest.ExpectedMap[k]
				if !ok {
					t.Errorf("failed to find '%s' from the generated map in the expected map", k)
				}

				if v != expVal {
					t.Errorf("value stored for '%s' (%+v) in the generated map not the same as what is in the expected map (%+v)", k, v, expVal)
				}
			}

			//compare expected to generated
			for k, v := range curTest.ExpectedMap {
				genVal, ok := genMap[k]
				if !ok {
					t.Errorf("failed to find '%s' from the expected map in the generated map", k)
				}

				if v != genVal {
					t.Errorf("value stored for '%s' (%+v) in the expedcted map not the same as what is in the generated map (%+v)", k, v, genVal)
				}
			}
		})
	}
}
