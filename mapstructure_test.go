package mapstructure

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Basic struct {
	VString     string
	VInt        int
	VUint       uint
	VBool       bool
	VFloat      float64
	VExtra      string
	vSilent     bool
	VData       interface{}
	VJsonInt    int
	VJsonFloat  float64
	VJsonNumber json.Number
}

type BasicPointer struct {
	VString     *string
	VInt        *int
	VUint       *uint
	VBool       *bool
	VFloat      *float64
	VExtra      *string
	vSilent     *bool
	VData       *interface{}
	VJsonInt    *int
	VJsonFloat  *float64
	VJsonNumber *json.Number
}

type Embedded struct {
	Basic
	VUnique string
}

type EmbeddedPointer struct {
	*Basic
	VUnique string
}

type SliceAlias []string

type EmbeddedSlice struct {
	SliceAlias `json:"slice_alias"`
	VUnique    string
}

type ArrayAlias [2]string

type EmbeddedArray struct {
	ArrayAlias `json:"array_alias"`
	VUnique    string
}

type Map struct {
	VFoo   string
	VOther map[string]string
}

type MapOfStruct struct {
	Value map[string]Basic
}

type Nested struct {
	VFoo string
	VBar Basic
}

type NestedPointer struct {
	VFoo string
	VBar *Basic
}

type Slice struct {
	VFoo string
	VBar []string
}

type SliceOfAlias struct {
	VFoo string
	VBar SliceAlias
}

type SliceOfStruct struct {
	Value []Basic
}

type SlicePointer struct {
	VBar *[]string
}

type Array struct {
	VFoo string
	VBar [2]string
}

type ArrayOfStruct struct {
	Value [2]Basic
}

type Func struct {
	Foo func() string
}

type Tagged struct {
	Extra string `json:"bar,what,what"`
	Value string `json:"foo"`
}

type StructWithOmitEmpty struct {
	VisibleStringField string                 `json:"visible-string"`
	OmitStringField    string                 `json:"omittable-string,omitempty"`
	VisibleIntField    int                    `json:"visible-int"`
	OmitIntField       int                    `json:"omittable-int,omitempty"`
	VisibleBoolField   bool                   `json:"visible-bool"`
	OmitBoolField      bool                   `json:"omittable-bool,omitempty"`
	VisibleFloatField  float64                `json:"visible-float"`
	OmitFloatField     float64                `json:"omittable-float,omitempty"`
	VisibleSliceField  []interface{}          `json:"visible-slice"`
	OmitSliceField     []interface{}          `json:"omittable-slice,omitempty"`
	VisibleMapField    map[string]interface{} `json:"visible-map"`
	OmitMapField       map[string]interface{} `json:"omittable-map,omitempty"`
	NestedField        *Nested                `json:"visible-nested"`
	OmitNestedField    *Nested                `json:"omittable-nested,omitempty"`
}

type TypeConversionResult struct {
	IntToFloat         float32
	IntToUint          uint
	IntToBool          bool
	IntToString        string
	UintToInt          int
	UintToFloat        float32
	UintToBool         bool
	UintToString       string
	BoolToInt          int
	BoolToUint         uint
	BoolToFloat        float32
	BoolToString       string
	FloatToInt         int
	FloatToUint        uint
	FloatToBool        bool
	FloatToString      string
	SliceUint8ToString string
	StringToSliceUint8 []byte
	ArrayUint8ToString string
	StringToInt        int
	StringToUint       uint
	StringToBool       bool
	StringToFloat      float32
	StringToStrSlice   []string
	StringToIntSlice   []int
	StringToStrArray   [1]string
	StringToIntArray   [1]int
	SliceToMap         map[string]interface{}
	MapToSlice         []interface{}
	ArrayToMap         map[string]interface{}
	MapToArray         [1]interface{}
}

func TestBasicTypes(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vstring":     "foo",
		"vint":        42,
		"VUint":       42,
		"vbool":       true,
		"VFloat":      42.42,
		"vSilent":     true,
		"vdata":       42,
		"vjsonInt":    json.Number("1234"),
		"vjsonFloat":  json.Number("1234.5"),
		"vjsonNumber": json.Number("1234.5"),
	}

	var result Basic
	err := Decode(input, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.VString != "foo" {
		t.Errorf("vstring value should be 'foo': %#v", result.VString)
	}

	if result.VInt != 42 {
		t.Errorf("vint value should be 42: %#v", result.VInt)
	}

	if result.VUint != 42 {
		t.Errorf("vuint value should be 42: %#v", result.VUint)
	}

	if result.VBool != true {
		t.Errorf("vbool value should be true: %#v", result.VBool)
	}

	if result.VFloat != 42.42 {
		t.Errorf("vfloat value should be 42.42: %#v", result.VFloat)
	}

	if result.VExtra != "" {
		t.Errorf("vextra value should be empty: %#v", result.VExtra)
	}

	if result.vSilent != false {
		t.Error("vSilent should not be set, it is unexported")
	}

	if result.VData != 42 {
		t.Error("vdata should be valid")
	}

	if result.VJsonInt != 1234 {
		t.Errorf("vjsonint value should be 1234: %#v", result.VJsonInt)
	}

	if result.VJsonFloat != 1234.5 {
		t.Errorf("vjsonfloat value should be 1234.5: %#v", result.VJsonFloat)
	}

	if !reflect.DeepEqual(result.VJsonNumber, json.Number("1234.5")) {
		t.Errorf("vjsonnumber value should be '1234.5': %T, %#v", result.VJsonNumber, result.VJsonNumber)
	}
}

func TestBasic_IntWithFloat(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vint": float64(42),
	}

	var result Basic
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err)
	}
}

func TestBasic_Merge(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vint": 42,
	}

	var result Basic
	result.VUint = 100
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err)
	}

	expected := Basic{
		VInt:  42,
		VUint: 100,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}

// Test for issue #46.
func TestBasic_Struct(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vdata": map[string]interface{}{
			"vstring": "foo",
		},
	}

	var result, inner Basic
	result.VData = &inner
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err)
	}
	expected := Basic{
		VData: &Basic{
			VString: "foo",
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestDecode_Embedded(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vstring": "foo",
		"Basic": map[string]interface{}{
			"vstring": "innerfoo",
		},
		"vunique": "bar",
	}

	var result Embedded
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if result.VString != "innerfoo" {
		t.Errorf("vstring value should be 'innerfoo': %#v", result.VString)
	}

	if result.VUnique != "bar" {
		t.Errorf("vunique value should be 'bar': %#v", result.VUnique)
	}
}

func TestDecode_EmbeddedPointer(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vstring": "foo",
		"Basic": map[string]interface{}{
			"vstring": "innerfoo",
		},
		"vunique": "bar",
	}

	var result EmbeddedPointer
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := EmbeddedPointer{
		Basic: &Basic{
			VString: "innerfoo",
		},
		VUnique: "bar",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestDecode_EmbeddedSlice(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"slice_alias": []string{"foo", "bar"},
		"vunique":     "bar",
	}

	var result EmbeddedSlice
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(result.SliceAlias, SliceAlias([]string{"foo", "bar"})) {
		t.Errorf("slice value: %#v", result.SliceAlias)
	}

	if result.VUnique != "bar" {
		t.Errorf("vunique value should be 'bar': %#v", result.VUnique)
	}
}

func TestDecode_EmbeddedArray(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"array_alias": [2]string{"foo", "bar"},
		"vunique":     "bar",
	}

	var result EmbeddedArray
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(result.ArrayAlias, ArrayAlias([2]string{"foo", "bar"})) {
		t.Errorf("array value: %#v", result.ArrayAlias)
	}

	if result.VUnique != "bar" {
		t.Errorf("vunique value should be 'bar': %#v", result.VUnique)
	}
}

func TestDecode_Nil(t *testing.T) {
	t.Parallel()

	var input interface{}
	result := Basic{
		VString: "foo",
	}

	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if result.VString != "foo" {
		t.Fatalf("bad: %#v", result.VString)
	}
}

func TestDecode_NonStruct(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	var result map[string]string
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if result["foo"] != "bar" {
		t.Fatal("foo is not bar")
	}
}

func TestDecode_StructMatch(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vbar": Basic{
			VString: "foo",
		},
	}

	var result Nested
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if result.VBar.VString != "foo" {
		t.Errorf("bad: %#v", result)
	}
}

func TestDecode_TypeConversion(t *testing.T) {
	input := map[string]interface{}{
		"IntToFloat":         42,
		"IntToUint":          42,
		"IntToBool":          1,
		"IntToString":        42,
		"UintToInt":          42,
		"UintToFloat":        42,
		"UintToBool":         42,
		"UintToString":       42,
		"BoolToInt":          true,
		"BoolToUint":         true,
		"BoolToFloat":        true,
		"BoolToString":       true,
		"FloatToInt":         42.42,
		"FloatToUint":        42.42,
		"FloatToBool":        42.42,
		"FloatToString":      42.42,
		"SliceUint8ToString": []uint8("foo"),
		"StringToSliceUint8": "foo",
		"ArrayUint8ToString": [3]uint8{'f', 'o', 'o'},
		"StringToInt":        "42",
		"StringToUint":       "42",
		"StringToBool":       "1",
		"StringToFloat":      "42.42",
		"StringToStrSlice":   "A",
		"StringToIntSlice":   "42",
		"StringToStrArray":   "A",
		"StringToIntArray":   "42",
		"SliceToMap":         []interface{}{},
		"MapToSlice":         map[string]interface{}{},
		"ArrayToMap":         []interface{}{},
		"MapToArray":         map[string]interface{}{},
	}

	expectedResultStrict := TypeConversionResult{
		IntToFloat:  42.0,
		IntToUint:   42,
		UintToInt:   42,
		UintToFloat: 42,
		BoolToInt:   0,
		BoolToUint:  0,
		BoolToFloat: 0,
		FloatToInt:  42,
		FloatToUint: 42,
	}

	// Test strict type conversion
	var resultStrict TypeConversionResult
	err := Decode(input, &resultStrict)
	if err == nil {
		t.Errorf("should return an error")
	}
	if !reflect.DeepEqual(resultStrict, expectedResultStrict) {
		t.Errorf("expected %v, got: %v", expectedResultStrict, resultStrict)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[interface{}]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}

	var result Map
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	if result.VFoo != "foo" {
		t.Errorf("vfoo value should be 'foo': %#v", result.VFoo)
	}

	if result.VOther == nil {
		t.Fatal("vother should not be nil")
	}

	if len(result.VOther) != 2 {
		t.Error("vother should have two items")
	}

	if result.VOther["foo"] != "foo" {
		t.Errorf("'foo' key should be foo, got: %#v", result.VOther["foo"])
	}

	if result.VOther["bar"] != "bar" {
		t.Errorf("'bar' key should be bar, got: %#v", result.VOther["bar"])
	}
}

func TestMapMerge(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[interface{}]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}

	var result Map
	result.VOther = map[string]string{"hello": "world"}
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an error: %s", err)
	}

	if result.VFoo != "foo" {
		t.Errorf("vfoo value should be 'foo': %#v", result.VFoo)
	}

	expected := map[string]string{
		"foo":   "foo",
		"bar":   "bar",
		"hello": "world",
	}
	if !reflect.DeepEqual(result.VOther, expected) {
		t.Errorf("bad: %#v", result.VOther)
	}
}

func TestMapOfStruct(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"value": map[string]interface{}{
			"foo": map[string]string{"vstring": "one"},
			"bar": map[string]string{"vstring": "two"},
		},
	}

	var result MapOfStruct
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err)
	}

	if result.Value == nil {
		t.Fatal("value should not be nil")
	}

	if len(result.Value) != 2 {
		t.Error("value should have two items")
	}

	if result.Value["foo"].VString != "one" {
		t.Errorf("foo value should be 'one', got: %s", result.Value["foo"].VString)
	}

	if result.Value["bar"].VString != "two" {
		t.Errorf("bar value should be 'two', got: %s", result.Value["bar"].VString)
	}
}

func TestNestedType(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": map[string]interface{}{
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		},
	}

	var result Nested
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if result.VFoo != "foo" {
		t.Errorf("vfoo value should be 'foo': %#v", result.VFoo)
	}

	if result.VBar.VString != "foo" {
		t.Errorf("vstring value should be 'foo': %#v", result.VBar.VString)
	}

	if result.VBar.VInt != 42 {
		t.Errorf("vint value should be 42: %#v", result.VBar.VInt)
	}

	if result.VBar.VBool != true {
		t.Errorf("vbool value should be true: %#v", result.VBar.VBool)
	}

	if result.VBar.VExtra != "" {
		t.Errorf("vextra value should be empty: %#v", result.VBar.VExtra)
	}
}

func TestNestedTypePointer(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": &map[string]interface{}{
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,
		},
	}

	var result NestedPointer
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if result.VFoo != "foo" {
		t.Errorf("vfoo value should be 'foo': %#v", result.VFoo)
	}

	if result.VBar.VString != "foo" {
		t.Errorf("vstring value should be 'foo': %#v", result.VBar.VString)
	}

	if result.VBar.VInt != 42 {
		t.Errorf("vint value should be 42: %#v", result.VBar.VInt)
	}

	if result.VBar.VBool != true {
		t.Errorf("vbool value should be true: %#v", result.VBar.VBool)
	}

	if result.VBar.VExtra != "" {
		t.Errorf("vextra value should be empty: %#v", result.VBar.VExtra)
	}
}

// Test for issue #46.
func TestNestedTypeInterface(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": &map[string]interface{}{
			"vstring": "foo",
			"vint":    42,
			"vbool":   true,

			"vdata": map[string]interface{}{
				"vstring": "bar",
			},
		},
	}

	var result NestedPointer
	result.VBar = new(Basic)
	result.VBar.VData = new(Basic)
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if result.VFoo != "foo" {
		t.Errorf("vfoo value should be 'foo': %#v", result.VFoo)
	}

	if result.VBar.VString != "foo" {
		t.Errorf("vstring value should be 'foo': %#v", result.VBar.VString)
	}

	if result.VBar.VInt != 42 {
		t.Errorf("vint value should be 42: %#v", result.VBar.VInt)
	}

	if result.VBar.VBool != true {
		t.Errorf("vbool value should be true: %#v", result.VBar.VBool)
	}

	if result.VBar.VExtra != "" {
		t.Errorf("vextra value should be empty: %#v", result.VBar.VExtra)
	}

	if result.VBar.VData.(*Basic).VString != "bar" {
		t.Errorf("vstring value should be 'bar': %#v", result.VBar.VData.(*Basic).VString)
	}
}

func TestSlice(t *testing.T) {
	t.Parallel()

	inputStringSlice := map[string]interface{}{
		"vfoo": "foo",
		"vbar": []string{"foo", "bar", "baz"},
	}

	inputStringSlicePointer := map[string]interface{}{
		"vfoo": "foo",
		"vbar": &[]string{"foo", "bar", "baz"},
	}

	outputStringSlice := &Slice{
		"foo",
		[]string{"foo", "bar", "baz"},
	}

	testSliceInput(t, inputStringSlice, outputStringSlice)
	testSliceInput(t, inputStringSlicePointer, outputStringSlice)
}

func TestInvalidSlice(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": 42,
	}

	result := Slice{}
	err := Decode(input, &result)
	if err == nil {
		t.Errorf("expected failure")
	}
}

func TestSliceOfStruct(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"value": []map[string]interface{}{
			{"vstring": "one"},
			{"vstring": "two"},
		},
	}

	var result SliceOfStruct
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	if len(result.Value) != 2 {
		t.Fatalf("expected two values, got %d", len(result.Value))
	}

	if result.Value[0].VString != "one" {
		t.Errorf("first value should be 'one', got: %s", result.Value[0].VString)
	}

	if result.Value[1].VString != "two" {
		t.Errorf("second value should be 'two', got: %s", result.Value[1].VString)
	}
}

func TestArray(t *testing.T) {
	t.Parallel()

	inputStringArray := map[string]interface{}{
		"vfoo": "foo",
		"vbar": [2]string{"foo", "bar"},
	}

	inputStringArrayPointer := map[string]interface{}{
		"vfoo": "foo",
		"vbar": &[2]string{"foo", "bar"},
	}

	outputStringArray := &Array{
		"foo",
		[2]string{"foo", "bar"},
	}

	testArrayInput(t, inputStringArray, outputStringArray)
	testArrayInput(t, inputStringArrayPointer, outputStringArray)
}

func TestInvalidArray(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": 42,
	}

	result := Array{}
	err := Decode(input, &result)
	if err == nil {
		t.Errorf("expected failure")
	}
}

func TestArrayOfStruct(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"value": []map[string]interface{}{
			{"vstring": "one"},
			{"vstring": "two"},
		},
	}

	var result ArrayOfStruct
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	if len(result.Value) != 2 {
		t.Fatalf("expected two values, got %d", len(result.Value))
	}

	if result.Value[0].VString != "one" {
		t.Errorf("first value should be 'one', got: %s", result.Value[0].VString)
	}

	if result.Value[1].VString != "two" {
		t.Errorf("second value should be 'two', got: %s", result.Value[1].VString)
	}
}

func TestDecodeTable(t *testing.T) {
	t.Parallel()

	// We need to make new types so that we don't get the short-circuit
	// copy functionality. We want to test the deep copying functionality.
	type BasicCopy Basic
	type NestedPointerCopy NestedPointer
	type MapCopy Map

	tests := []struct {
		name    string
		in      interface{}
		target  interface{}
		out     interface{}
		wantErr bool
	}{
		{
			"basic struct input",
			&Basic{
				VString: "vstring",
				VInt:    2,
				VUint:   3,
				VBool:   true,
				VFloat:  4.56,
				VExtra:  "vextra",
				vSilent: true,
				VData:   []byte("data"),
			},
			&map[string]interface{}{},
			&map[string]interface{}{
				"VString":     "vstring",
				"VInt":        2,
				"VUint":       uint(3),
				"VBool":       true,
				"VFloat":      4.56,
				"VExtra":      "vextra",
				"VData":       []byte("data"),
				"VJsonInt":    0,
				"VJsonFloat":  0.0,
				"VJsonNumber": json.Number(""),
			},
			false,
		},
		{
			"embedded struct input",
			&Embedded{
				VUnique: "vunique",
				Basic: Basic{
					VString: "vstring",
					VInt:    2,
					VUint:   3,
					VBool:   true,
					VFloat:  4.56,
					VExtra:  "vextra",
					vSilent: true,
					VData:   []byte("data"),
				},
			},
			&map[string]interface{}{},
			&map[string]interface{}{
				"VUnique": "vunique",
				"Basic": map[string]interface{}{
					"VString":     "vstring",
					"VInt":        2,
					"VUint":       uint(3),
					"VBool":       true,
					"VFloat":      4.56,
					"VExtra":      "vextra",
					"VData":       []byte("data"),
					"VJsonInt":    0,
					"VJsonFloat":  0.0,
					"VJsonNumber": json.Number(""),
				},
			},
			false,
		},
		{
			"struct => struct",
			&Basic{
				VString: "vstring",
				VInt:    2,
				VUint:   3,
				VBool:   true,
				VFloat:  4.56,
				VExtra:  "vextra",
				VData:   []byte("data"),
				vSilent: true,
			},
			&BasicCopy{},
			&BasicCopy{
				VString: "vstring",
				VInt:    2,
				VUint:   3,
				VBool:   true,
				VFloat:  4.56,
				VExtra:  "vextra",
				VData:   []byte("data"),
			},
			false,
		},
		{
			"struct => struct with pointers",
			&NestedPointer{
				VFoo: "hello",
				VBar: nil,
			},
			&NestedPointerCopy{},
			&NestedPointerCopy{
				VFoo: "hello",
			},
			false,
		},
		{
			"basic pointer to non-pointer",
			&BasicPointer{
				VString: stringPtr("vstring"),
				VInt:    intPtr(2),
				VUint:   uintPtr(3),
				VBool:   boolPtr(true),
				VFloat:  floatPtr(4.56),
				VData:   interfacePtr([]byte("data")),
			},
			&Basic{},
			&Basic{
				VString: "vstring",
				VInt:    2,
				VUint:   3,
				VBool:   true,
				VFloat:  4.56,
				VData:   []byte("data"),
			},
			false,
		},
		{
			"slice non-pointer to pointer",
			&Slice{},
			&SlicePointer{},
			&SlicePointer{},
			false,
		},
		{
			"slice non-pointer to pointer, zero field",
			&Slice{},
			&SlicePointer{
				VBar: &[]string{"yo"},
			},
			&SlicePointer{},
			false,
		},
		{
			"slice to slice alias",
			&Slice{},
			&SliceOfAlias{},
			&SliceOfAlias{},
			false,
		},
		{
			"nil map to map",
			&Map{},
			&MapCopy{},
			&MapCopy{},
			false,
		},
		{
			"nil map to non-empty map",
			&Map{},
			&MapCopy{VOther: map[string]string{"foo": "bar"}},
			&MapCopy{},
			false,
		},

		{
			"slice input - should error",
			[]string{"foo", "bar"},
			&map[string]interface{}{},
			&map[string]interface{}{},
			true,
		},
		{
			"struct with slice property",
			&Slice{
				VFoo: "vfoo",
				VBar: []string{"foo", "bar"},
			},
			&map[string]interface{}{},
			&map[string]interface{}{
				"VFoo": "vfoo",
				"VBar": []string{"foo", "bar"},
			},
			false,
		},
		{
			"struct with slice of struct property",
			&SliceOfStruct{
				Value: []Basic{
					Basic{
						VString: "vstring",
						VInt:    2,
						VUint:   3,
						VBool:   true,
						VFloat:  4.56,
						VExtra:  "vextra",
						vSilent: true,
						VData:   []byte("data"),
					},
				},
			},
			&map[string]interface{}{},
			&map[string]interface{}{
				"Value": []Basic{
					Basic{
						VString: "vstring",
						VInt:    2,
						VUint:   3,
						VBool:   true,
						VFloat:  4.56,
						VExtra:  "vextra",
						vSilent: true,
						VData:   []byte("data"),
					},
				},
			},
			false,
		},
		{
			"struct with map property",
			&Map{
				VFoo:   "vfoo",
				VOther: map[string]string{"vother": "vother"},
			},
			&map[string]interface{}{},
			&map[string]interface{}{
				"VFoo": "vfoo",
				"VOther": map[string]string{
					"vother": "vother",
				}},
			false,
		},
		{
			"tagged struct",
			&Tagged{
				Extra: "extra",
				Value: "value",
			},
			&map[string]string{},
			&map[string]string{
				"bar": "extra",
				"foo": "value",
			},
			false,
		},
		{
			"omit tag struct",
			&struct {
				Value string `json:"value"`
				Omit  string `json:"-"`
			}{
				Value: "value",
				Omit:  "omit",
			},
			&map[string]string{},
			&map[string]string{
				"value": "value",
			},
			false,
		},
		{
			"decode to wrong map type",
			&struct {
				Value string
			}{
				Value: "string",
			},
			&map[string]int{},
			&map[string]int{},
			true,
		},
		{
			"struct with omitempty tag return non-empty values",
			&struct {
				VisibleField interface{} `json:"visible"`
				OmitField    interface{} `json:"omittable,omitempty"`
			}{
				VisibleField: nil,
				OmitField:    "string",
			},
			&map[string]interface{}{},
			&map[string]interface{}{"visible": nil, "omittable": "string"},
			false,
		},
		{
			"struct with omitempty tag ignore empty values",
			&struct {
				VisibleField interface{} `json:"visible"`
				OmitField    interface{} `json:"omittable,omitempty"`
			}{
				VisibleField: nil,
				OmitField:    nil,
			},
			&map[string]interface{}{},
			&map[string]interface{}{"visible": nil},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.in, tt.target); (err != nil) != tt.wantErr {
				t.Fatalf("%q: TestMapOutputForStructuredInputs() unexpected error: %s", tt.name, err)
			}

			if !reflect.DeepEqual(tt.out, tt.target) {
				t.Fatalf("%q: TestMapOutputForStructuredInputs() expected: %#v, got: %#v", tt.name, tt.out, tt.target)
			}
		})
	}
}

func TestInvalidType(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"vstring": 42,
	}

	var result Basic
	err := Decode(input, &result)
	if err == nil {
		t.Fatal("error should exist")
	}

	derr, ok := err.(*Error)
	if !ok {
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	}

	if derr.Errors[0] != "'VString' expected type 'string', got unconvertible type 'int'" {
		t.Errorf("got unexpected error: %s", err)
	}

	inputNegIntUint := map[string]interface{}{
		"vuint": -42,
	}

	err = Decode(inputNegIntUint, &result)
	if err == nil {
		t.Fatal("error should exist")
	}

	derr, ok = err.(*Error)
	if !ok {
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	}

	if derr.Errors[0] != "cannot parse 'VUint', -42 overflows uint" {
		t.Errorf("got unexpected error: %s", err)
	}

	inputNegFloatUint := map[string]interface{}{
		"vuint": -42.0,
	}

	err = Decode(inputNegFloatUint, &result)
	if err == nil {
		t.Fatal("error should exist")
	}

	derr, ok = err.(*Error)
	if !ok {
		t.Fatalf("error should be kind of Error, instead: %#v", err)
	}

	if derr.Errors[0] != "cannot parse 'VUint', -42.000000 overflows uint" {
		t.Errorf("got unexpected error: %s", err)
	}
}

func TestNonPtrValue(t *testing.T) {
	t.Parallel()

	err := Decode(map[string]interface{}{}, Basic{})
	if err == nil {
		t.Fatal("error should exist")
	}

	if err.Error() != "result must be a pointer" {
		t.Errorf("got unexpected error: %s", err)
	}
}

func TestTagged(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"foo": "bar",
		"bar": "value",
	}

	var result Tagged
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if result.Value != "bar" {
		t.Errorf("value should be 'bar', got: %#v", result.Value)
	}

	if result.Extra != "value" {
		t.Errorf("extra should be 'value', got: %#v", result.Extra)
	}
}

func TestDecode_StructTaggedWithOmitempty_OmitEmptyValues(t *testing.T) {
	t.Parallel()

	input := &StructWithOmitEmpty{}

	var emptySlice []interface{}
	var emptyMap map[string]interface{}
	var emptyNested *Nested
	expected := &map[string]interface{}{
		"visible-string": "",
		"visible-int":    0,
		"visible-bool":   false,
		"visible-float":  0.0,
		"visible-slice":  emptySlice,
		"visible-map":    emptyMap,
		"visible-nested": emptyNested,
	}

	actual := &map[string]interface{}{}
	Decode(input, actual)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Decode() expected: %#v, got: %#v", expected, actual)
	}
}

func TestDecode_StructTaggedWithOmitempty_KeepNonEmptyValues(t *testing.T) {
	t.Parallel()

	input := &StructWithOmitEmpty{
		VisibleStringField: "",
		OmitStringField:    "string",
		VisibleIntField:    0,
		OmitIntField:       1,
		VisibleBoolField:   false,
		OmitBoolField:      true,
		VisibleFloatField:  0.0,
		OmitFloatField:     1.0,
		VisibleSliceField:  nil,
		OmitSliceField:     []interface{}{1},
		VisibleMapField:    nil,
		OmitMapField:       map[string]interface{}{"k": "v"},
		NestedField:        nil,
		OmitNestedField:    &Nested{},
	}

	var emptySlice []interface{}
	var emptyMap map[string]interface{}
	var emptyNested *Nested
	expected := &map[string]interface{}{
		"visible-string":   "",
		"omittable-string": "string",
		"visible-int":      0,
		"omittable-int":    1,
		"visible-bool":     false,
		"omittable-bool":   true,
		"visible-float":    0.0,
		"omittable-float":  1.0,
		"visible-slice":    emptySlice,
		"omittable-slice":  []interface{}{1},
		"visible-map":      emptyMap,
		"omittable-map":    map[string]interface{}{"k": "v"},
		"visible-nested":   emptyNested,
		"omittable-nested": &Nested{},
	}

	actual := &map[string]interface{}{}
	Decode(input, actual)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Decode() expected: %#v, got: %#v", expected, actual)
	}
}

func testSliceInput(t *testing.T, input map[string]interface{}, expected *Slice) {
	var result Slice
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if result.VFoo != expected.VFoo {
		t.Errorf("VFoo expected '%s', got '%s'", expected.VFoo, result.VFoo)
	}

	if result.VBar == nil {
		t.Fatalf("VBar a slice, got '%#v'", result.VBar)
	}

	if len(result.VBar) != len(expected.VBar) {
		t.Errorf("VBar length should be %d, got %d", len(expected.VBar), len(result.VBar))
	}

	for i, v := range result.VBar {
		if v != expected.VBar[i] {
			t.Errorf(
				"VBar[%d] should be '%#v', got '%#v'",
				i, expected.VBar[i], v)
		}
	}
}

func testArrayInput(t *testing.T, input map[string]interface{}, expected *Array) {
	var result Array
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if result.VFoo != expected.VFoo {
		t.Errorf("VFoo expected '%s', got '%s'", expected.VFoo, result.VFoo)
	}

	if result.VBar == [2]string{} {
		t.Fatalf("VBar a slice, got '%#v'", result.VBar)
	}

	if len(result.VBar) != len(expected.VBar) {
		t.Errorf("VBar length should be %d, got %d", len(expected.VBar), len(result.VBar))
	}

	for i, v := range result.VBar {
		if v != expected.VBar[i] {
			t.Errorf(
				"VBar[%d] should be '%#v', got '%#v'",
				i, expected.VBar[i], v)
		}
	}
}

func TestDecode_Func(t *testing.T) {
	input := map[string]interface{}{
		"foo": func() string { return "baz" },
	}

	var result Func
	err := Decode(input, &result)
	if err != nil {
		t.Fatalf("got an err: %s", err)
	}

	if result.Foo() != "baz" {
		t.Errorf("Foo call result should be 'baz': %s", result.Foo())
	}
}

func stringPtr(v string) *string              { return &v }
func intPtr(v int) *int                       { return &v }
func uintPtr(v uint) *uint                    { return &v }
func boolPtr(v bool) *bool                    { return &v }
func floatPtr(v float64) *float64             { return &v }
func interfacePtr(v interface{}) *interface{} { return &v }
