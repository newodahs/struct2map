package struct2map

import (
	"testing"
)

type simpleTestStruct struct {
	notExportedField           int `struct2map:"notExported"`
	NegatedExportField         int `struct2map:"-"`
	RegularFieldNoTag          int
	RegularFieldNameTag        int   `struct2map:"regularField"`
	RegularFieldOmitEmpty      *int  `struct2map:"regularFieldOmitted,omitempty"`
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
}

type embeddedStruct struct {
	TopLevelValue              bool             `struct2map:"topLevelBool"`
	notExportedStruct          simpleTestStruct `struct2map:"notExportedStruct"`
	NegatedExportedStruct      simpleTestStruct `struct2map:"-"`
	RegularExportStructNoTag   simpleTestStruct
	RegularFieldNameTag        simpleTestStruct   `struct2map:"regularStruct"`
	RegularFieldOmitEmpty      *simpleTestStruct  `struct2map:"regularStructOmitted,omitempty"`
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

	testStructPtr := &simpleTestStruct{
		notExportedField:           simpleInt,
		NegatedExportField:         simpleInt,
		RegularFieldNoTag:          simpleInt,
		RegularFieldNameTag:        simpleInt,
		RegularFieldOmitEmpty:      &simpleInt,
		RegularFieldPointerPointer: &simpleIntPtr,
	}

	testSet := []struct {
		Name          string
		TestStructure any
		ExpectedMap   map[string]any
		SkipTest      bool
	}{
		{
			Name: "simpleTestStruct Validation",
			TestStructure: simpleTestStruct{
				notExportedField:           simpleInt,
				NegatedExportField:         simpleInt,
				RegularFieldNoTag:          simpleInt,
				RegularFieldNameTag:        simpleInt,
				RegularFieldOmitEmpty:      &simpleInt,
				RegularFieldPointerPointer: &simpleIntPtr,
			},
			ExpectedMap: map[string]any{"RegularFieldNoTag": 1, "regularField": 1, "regularFieldOmitted": simpleInt, "regularFieldPointerPointer": simpleInt},
		},
		{
			Name: "pointer of simpleTestStruct Validation",
			TestStructure: &simpleTestStruct{
				notExportedField:           simpleInt,
				NegatedExportField:         simpleInt,
				RegularFieldNoTag:          simpleInt,
				RegularFieldNameTag:        simpleInt,
				RegularFieldOmitEmpty:      &simpleInt,
				RegularFieldPointerPointer: &simpleIntPtr,
			},
			ExpectedMap: map[string]any{"RegularFieldNoTag": 1, "regularField": 1, "regularFieldOmitted": simpleInt, "regularFieldPointerPointer": simpleInt},
		},
		{
			Name: "simpleTestStruct Validation: RegularFieldOmitEmpty nil to omit; RegularFieldPointerPointer nil and not omitted",
			TestStructure: simpleTestStruct{
				notExportedField:    simpleInt,
				NegatedExportField:  simpleInt,
				RegularFieldNoTag:   simpleInt,
				RegularFieldNameTag: simpleInt,
			},
			ExpectedMap: map[string]any{"RegularFieldNoTag": simpleInt, "regularField": simpleInt, "regularFieldPointerPointer": nil},
		},
		{
			Name: "complexTestStruct Validation",
			TestStructure: complexTestStruct{
				TopLevelField:              true,
				SliceField:                 []int{simpleInt, simpleInt, simpleInt},
				SliceFieldPtrVal:           []*int{&simpleInt, &simpleInt, &simpleInt},
				MapFieldStrKey:             map[string]int{"field-a": simpleInt, "field-b": simpleInt},
				MapFieldStrKeyPtrVal:       map[string]*int{"field-a-ptr": &simpleInt, "field-b-ptr": &simpleInt},
				MapFieldIntKey:             map[int]string{1: "test1", 2: "test2"},
				MapFieldStrKeyStructVal:    map[string]simpleTestStruct{"simpleStruct1": *testStructPtr},
				MapFieldStrKeyStructPtrVal: map[string]*simpleTestStruct{"simpleStructPtr1": testStructPtr},
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
				"mapFieldStrKeyStructPtrVal.regularFieldOmitted":        1,
				"mapFieldStrKeyStructPtrVal.regularFieldPointerPointer": 1,
				"mapFieldStrKeyStructVal.RegularFieldNoTag":             1,
				"mapFieldStrKeyStructVal.regularField":                  1,
				"mapFieldStrKeyStructVal.regularFieldOmitted":           1,
				"mapFieldStrKeyStructVal.regularFieldPointerPointer":    1,
				"sliceField.0":                                          1,
				"sliceField.1":                                          1,
				"sliceField.2":                                          1,
				"sliceFieldPtrVal.0":                                    1,
				"sliceFieldPtrVal.1":                                    1,
				"sliceFieldPtrVal.2":                                    1,
				"topLevelBool":                                          true,
			},
		},
		{
			Name: "embeddedStruct Pointer Validation",
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
					RegStruct:    *testStructPtr,
					RegStructPtr: testStructPtr,
				},
			},
			ExpectedMap: map[string]any{
				"RegularExportStructNoTag.RegularFieldNoTag":             1,
				"RegularExportStructNoTag.regularField":                  1,
				"RegularExportStructNoTag.regularFieldOmitted":           1,
				"RegularExportStructNoTag.regularFieldPointerPointer":    1,
				"anonStruct.RegStruct.RegularFieldNoTag":                 1,
				"anonStruct.RegStruct.regularField":                      1,
				"anonStruct.RegStruct.regularFieldOmitted":               1,
				"anonStruct.RegStruct.regularFieldPointerPointer":        1,
				"anonStruct.RegStructPtr.RegularFieldNoTag":              1,
				"anonStruct.RegStructPtr.regularField":                   1,
				"anonStruct.RegStructPtr.regularFieldOmitted":            1,
				"anonStruct.RegStructPtr.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStruct.RegularFieldNoTag":              1,
				"anonStructPtr.RegStruct.regularField":                   1,
				"anonStructPtr.RegStruct.regularFieldOmitted":            1,
				"anonStructPtr.RegStruct.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStructPtr.RegularFieldNoTag":           1,
				"anonStructPtr.RegStructPtr.regularField":                1,
				"anonStructPtr.RegStructPtr.regularFieldOmitted":         1,
				"anonStructPtr.RegStructPtr.regularFieldPointerPointer":  1,
				"regularStruct.RegularFieldNoTag":                        1,
				"regularStruct.regularField":                             1,
				"regularStruct.regularFieldOmitted":                      1,
				"regularStruct.regularFieldPointerPointer":               1,
				"regularStructOmitted.RegularFieldNoTag":                 1,
				"regularStructOmitted.regularField":                      1,
				"regularStructOmitted.regularFieldOmitted":               1,
				"regularStructOmitted.regularFieldPointerPointer":        1,
				"regularStructPointerPointer.RegularFieldNoTag":          1,
				"regularStructPointerPointer.regularField":               1,
				"regularStructPointerPointer.regularFieldOmitted":        1,
				"regularStructPointerPointer.regularFieldPointerPointer": 1,
				"topLevelBool": true,
			},
		},
		{
			Name: "pointer of embeddedStruct Validation",
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
					RegStruct:    *testStructPtr,
					RegStructPtr: testStructPtr,
				},
				AnonStructPtr: &struct {
					RegStruct    simpleTestStruct
					RegStructPtr *simpleTestStruct
				}{
					RegStruct:    *testStructPtr,
					RegStructPtr: testStructPtr,
				},
			},
			ExpectedMap: map[string]any{
				"RegularExportStructNoTag.RegularFieldNoTag":             1,
				"RegularExportStructNoTag.regularField":                  1,
				"RegularExportStructNoTag.regularFieldOmitted":           1,
				"RegularExportStructNoTag.regularFieldPointerPointer":    1,
				"anonStruct.RegStruct.RegularFieldNoTag":                 1,
				"anonStruct.RegStruct.regularField":                      1,
				"anonStruct.RegStruct.regularFieldOmitted":               1,
				"anonStruct.RegStruct.regularFieldPointerPointer":        1,
				"anonStruct.RegStructPtr.RegularFieldNoTag":              1,
				"anonStruct.RegStructPtr.regularField":                   1,
				"anonStruct.RegStructPtr.regularFieldOmitted":            1,
				"anonStruct.RegStructPtr.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStruct.RegularFieldNoTag":              1,
				"anonStructPtr.RegStruct.regularField":                   1,
				"anonStructPtr.RegStruct.regularFieldOmitted":            1,
				"anonStructPtr.RegStruct.regularFieldPointerPointer":     1,
				"anonStructPtr.RegStructPtr.RegularFieldNoTag":           1,
				"anonStructPtr.RegStructPtr.regularField":                1,
				"anonStructPtr.RegStructPtr.regularFieldOmitted":         1,
				"anonStructPtr.RegStructPtr.regularFieldPointerPointer":  1,
				"regularStruct.RegularFieldNoTag":                        1,
				"regularStruct.regularField":                             1,
				"regularStruct.regularFieldOmitted":                      1,
				"regularStruct.regularFieldPointerPointer":               1,
				"regularStructOmitted.RegularFieldNoTag":                 1,
				"regularStructOmitted.regularField":                      1,
				"regularStructOmitted.regularFieldOmitted":               1,
				"regularStructOmitted.regularFieldPointerPointer":        1,
				"regularStructPointerPointer.RegularFieldNoTag":          1,
				"regularStructPointerPointer.regularField":               1,
				"regularStructPointerPointer.regularFieldOmitted":        1,
				"regularStructPointerPointer.regularFieldPointerPointer": 1,
				"topLevelBool": true,
			},
		},
		{
			Name: "flattenStruct Validation: anonContained should not appear as a parent",
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
				"structIgnoreParent.regularFieldOmitted":        1,
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
