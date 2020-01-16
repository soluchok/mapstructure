package mapstructure

import (
	"encoding/json"
	"testing"
)

type Person struct {
	Name   string
	Age    int
	Emails []string
	Extra  map[string]string
}

func Benchmark_Decode(b *testing.B) {
	input := map[string]interface{}{
		"name":   "Benchmark",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"nested": "benchmark",
		},
	}

	var result Person
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

// decodeViaJSON takes the map data and passes it through encoding/json to convert it into the
// given Go native structure pointed to by v. v must be a pointer to a struct.
func decodeViaJSON(data interface{}, v interface{}) error {
	// Perform the task by simply marshalling the input into JSON,
	// then unmarshalling it into target native Go struct.
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func Benchmark_DecodeViaJSON(b *testing.B) {
	input := map[string]interface{}{
		"name":   "Benchmark",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"nested": "benchmark",
		},
	}

	var result Person
	for i := 0; i < b.N; i++ {
		decodeViaJSON(input, &result)
	}
}

func Benchmark_JSONUnmarshal(b *testing.B) {
	input := map[string]interface{}{
		"name":   "Benchmark",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"nested": "benchmark",
		},
	}

	inputB, err := json.Marshal(input)
	if err != nil {
		b.Fatal("Failed to marshal test input:", err)
	}

	var result Person
	for i := 0; i < b.N; i++ {
		json.Unmarshal(inputB, &result)
	}
}

func Benchmark_DecodeBasic(b *testing.B) {
	input := map[string]interface{}{
		"vstring":     "foo",
		"vint":        42,
		"Vuint":       42,
		"vbool":       true,
		"Vfloat":      42.42,
		"vsilent":     true,
		"vdata":       42,
		"vjsonInt":    json.Number("1234"),
		"vjsonFloat":  json.Number("1234.5"),
		"vjsonNumber": json.Number("1234.5"),
	}

	for i := 0; i < b.N; i++ {
		var result Basic
		Decode(input, &result)
	}
}

func Benchmark_DecodeEmbedded(b *testing.B) {
	input := map[string]interface{}{
		"vstring": "foo",
		"Basic": map[string]interface{}{
			"vstring": "innerfoo",
		},
		"vunique": "bar",
	}

	var result Embedded
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

func Benchmark_DecodeTypeConversion(b *testing.B) {
	input := map[string]interface{}{
		"IntToFloat":    42,
		"IntToUint":     42,
		"IntToBool":     1,
		"IntToString":   42,
		"UintToInt":     42,
		"UintToFloat":   42,
		"UintToBool":    42,
		"UintToString":  42,
		"BoolToInt":     true,
		"BoolToUint":    true,
		"BoolToFloat":   true,
		"BoolToString":  true,
		"FloatToInt":    42.42,
		"FloatToUint":   42.42,
		"FloatToBool":   42.42,
		"FloatToString": 42.42,
		"StringToInt":   "42",
		"StringToUint":  "42",
		"StringToBool":  "1",
		"StringToFloat": "42.42",
		"SliceToMap":    []interface{}{},
		"MapToSlice":    map[string]interface{}{},
	}

	var resultStrict TypeConversionResult
	for i := 0; i < b.N; i++ {
		Decode(input, &resultStrict)
	}
}

func Benchmark_DecodeMap(b *testing.B) {
	input := map[string]interface{}{
		"vfoo": "foo",
		"vother": map[interface{}]interface{}{
			"foo": "foo",
			"bar": "bar",
		},
	}

	var result Map
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

func Benchmark_DecodeMapOfStruct(b *testing.B) {
	input := map[string]interface{}{
		"value": map[string]interface{}{
			"foo": map[string]string{"vstring": "one"},
			"bar": map[string]string{"vstring": "two"},
		},
	}

	var result MapOfStruct
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

func Benchmark_DecodeSlice(b *testing.B) {
	input := map[string]interface{}{
		"vfoo": "foo",
		"vbar": []string{"foo", "bar", "baz"},
	}

	var result Slice
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

func Benchmark_DecodeSliceOfStruct(b *testing.B) {
	input := map[string]interface{}{
		"value": []map[string]interface{}{
			{"vstring": "one"},
			{"vstring": "two"},
		},
	}

	var result SliceOfStruct
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}

func Benchmark_DecodeTagged(b *testing.B) {
	input := map[string]interface{}{
		"foo": "bar",
		"bar": "value",
	}

	var result Tagged
	for i := 0; i < b.N; i++ {
		Decode(input, &result)
	}
}
