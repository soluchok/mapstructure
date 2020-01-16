// Package mapstructure exposes functionality to convert an arbitrary
// map[string]interface{} into a native Go structure.
//
// The Go structure can be arbitrarily complex, containing slices,
// other structs, etc. and the decoder will properly decode nested
// maps and so on into the proper structures in the native Go struct.
// See the examples to see what the decoder is capable of.
package mapstructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const tagName = `json`

// A Decoder takes a raw interface value and turns it into structured
// data, keeping track of rich error information along the way in case
// anything goes wrong.
type Decoder struct {
	// Result is a pointer to the struct that will contain the decoded
	// value.
	Result interface{}
}

// Decode takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a map or struct.
func Decode(input interface{}, output interface{}) error {
	val := reflect.ValueOf(output)
	if val.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	if !val.Elem().CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}

	return (&Decoder{Result: output}).Decode(input)
}

// Decode decodes the given raw interface to the target pointer specified
// by the configuration.
func (d *Decoder) Decode(input interface{}) error {
	return d.decode("", input, reflect.ValueOf(d.Result).Elem())
}

// Decodes an unknown data type into a specific reflection value.
func (d *Decoder) decode(name string, input interface{}, outVal reflect.Value) error {
	var inputVal reflect.Value
	if input != nil {
		inputVal = reflect.ValueOf(input)

		// We need to check here if input is a typed nil. Typed nils won't
		// match the "input == nil" below so we check that here.
		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if input == nil {
		return nil
	}

	if !inputVal.IsValid() {
		// If the input value is invalid, then we just set the value
		// to be the zero value.
		outVal.Set(reflect.Zero(outVal.Type()))
		return nil
	}

	switch getKind(outVal) {
	case reflect.Bool:
		return d.decodeBool(name, input, outVal)
	case reflect.Interface:
		return d.decodeBasic(name, input, outVal)
	case reflect.String:
		return d.decodeString(name, input, outVal)
	case reflect.Int:
		return d.decodeInt(name, input, outVal)
	case reflect.Uint:
		return d.decodeUint(name, input, outVal)
	case reflect.Float32:
		return d.decodeFloat(name, input, outVal)
	case reflect.Struct:
		return d.decodeStruct(name, input, outVal)
	case reflect.Map:
		return d.decodeMap(name, input, outVal)
	case reflect.Ptr:
		return d.decodePtr(name, input, outVal)
	case reflect.Slice:
		return d.decodeSlice(name, input, outVal)
	case reflect.Array:
		return d.decodeArray(name, input, outVal)
	case reflect.Func:
		return d.decodeFunc(name, input, outVal)
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("%s: unsupported type: %s", name, getKind(outVal))
	}
}

// This decodes a basic type (bool, int, string, etc.) and sets the
// value to "data" of that type.
func (d *Decoder) decodeBasic(name string, data interface{}, val reflect.Value) error {
	if val.IsValid() && val.Elem().IsValid() {
		return d.decode(name, data, val.Elem())
	}

	dataVal := reflect.ValueOf(data)

	// If the input data is a pointer, and the assigned type is the dereference
	// of that exact pointer, then indirect it so that we can assign it.
	// Example: *string to string
	if dataVal.Kind() == reflect.Ptr && dataVal.Type().Elem() == val.Type() {
		dataVal = reflect.Indirect(dataVal)
	}

	if !dataVal.IsValid() {
		dataVal = reflect.Zero(val.Type())
	}

	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) {
		return fmt.Errorf(
			"'%s' expected type '%s', got '%s'",
			name, val.Type(), dataValType)
	}

	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeString(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	if dataKind == reflect.String {
		val.SetString(dataVal.String())
		return nil
	}

	return fmt.Errorf("'%s' expected type '%s', got unconvertible type '%s'", name, val.Type(), dataVal.Type())
}

func (d *Decoder) decodeInt(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf("error decoding json.Number into %s: %s", name, err)
		}
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeUint(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Int:
		i := dataVal.Int()
		if i < 0 {
			return fmt.Errorf("cannot parse '%s', %d overflows uint",
				name, i)
		}
		val.SetUint(uint64(i))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		f := dataVal.Float()
		if f < 0 {
			return fmt.Errorf("cannot parse '%s', %f overflows uint",
				name, f)
		}
		val.SetUint(uint64(f))
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeBool(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	if dataKind != reflect.Bool {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'", name, val.Type(), dataVal.Type())
	}

	val.SetBool(dataVal.Bool())

	return nil
}

func (d *Decoder) decodeFloat(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(dataVal.Float())
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Float64()
		if err != nil {
			return fmt.Errorf("error decoding json.Number into %s: %s", name, err)
		}
		val.SetFloat(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeMap(name string, data interface{}, val reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// By default we overwrite keys in the current map
	valMap := val

	// If the map is nil or we're purposely zeroing fields, make a new map
	if valMap.IsNil() {
		// Make a new map to hold our result
		mapType := reflect.MapOf(valKeyType, valElemType)
		valMap = reflect.MakeMap(mapType)
	}

	// Check input type and based on the input type jump to the proper func
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	switch dataVal.Kind() {
	case reflect.Map:
		return d.decodeMapFromMap(name, dataVal, val, valMap)

	case reflect.Struct:
		return d.decodeMapFromStruct(name, dataVal, val, valMap)
	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeMapFromMap(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// Accumulate errors
	errors := make([]string, 0)

	// If the input data is empty, then we just match what the input data is.
	if dataVal.Len() == 0 {
		if dataVal.IsNil() {
			if !val.IsNil() {
				val.Set(dataVal)
			}
		} else {
			// Set to empty allocated value
			val.Set(valMap)
		}

		return nil
	}

	for _, k := range dataVal.MapKeys() {
		fieldName := fmt.Sprintf("%s[%s]", name, k)

		// First decode the key into the proper type
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := d.decode(fieldName, k.Interface(), currentKey); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		// Next decode the data into the proper type
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := d.decode(fieldName, v, currentVal); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		valMap.SetMapIndex(currentKey, currentVal)
	}

	// Set the built up map to the value
	val.Set(valMap)

	// If we had errors, return those
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeMapFromStruct(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	typ := dataVal.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Get the StructField first since this is a cheap operation. If the
		// field is unexported, then ignore it.
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		// Next get the actual value of this field and verify it is assignable
		// to the map value.
		v := dataVal.Field(i)
		if !v.Type().AssignableTo(valMap.Type().Elem()) {
			return fmt.Errorf("cannot assign type '%s' to map value field of type '%s'", v.Type(), valMap.Type().Elem())
		}

		tagValue := f.Tag.Get(tagName)
		tagParts := strings.Split(tagValue, ",")

		// If "omitempty" is specified in the tag, it ignores empty values.
		omitempty := false
		for _, tag := range tagParts[1:] {
			if tag == "omitempty" {
				omitempty = true
				break
			}
		}
		if omitempty && isEmptyValue(v) {
			continue
		}

		// Determine the name of the key in the map
		keyName := f.Name
		if tagParts[0] != "" {
			if tagParts[0] == "-" {
				continue
			}
			keyName = tagParts[0]
		}

		switch v.Kind() {
		// this is an embedded struct, so handle it differently
		case reflect.Struct:
			x := reflect.New(v.Type())
			x.Elem().Set(v)

			vType := valMap.Type()
			vKeyType := vType.Key()
			vElemType := vType.Elem()
			mType := reflect.MapOf(vKeyType, vElemType)
			vMap := reflect.MakeMap(mType)

			err := d.decode(keyName, x.Interface(), vMap)
			if err != nil {
				return err
			}
			valMap.SetMapIndex(reflect.ValueOf(keyName), vMap)
		default:
			valMap.SetMapIndex(reflect.ValueOf(keyName), v)
		}
	}

	if val.CanAddr() {
		val.Set(valMap)
	}

	return nil
}

func (d *Decoder) decodePtr(name string, data interface{}, val reflect.Value) error {
	// If the input data is nil, then we want to just set the output
	// pointer to be nil as well.
	isNil := data == nil
	if !isNil {
		switch v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			isNil = v.IsNil()
		}
	}
	if isNil {
		if !val.IsNil() && val.CanSet() {
			nilValue := reflect.New(val.Type()).Elem()
			val.Set(nilValue)
		}

		return nil
	}

	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	valType := val.Type()
	valElemType := valType.Elem()
	if val.CanSet() {
		realVal := val
		if realVal.IsNil() {
			realVal = reflect.New(valElemType)
		}

		if err := d.decode(name, data, reflect.Indirect(realVal)); err != nil {
			return err
		}

		val.Set(realVal)
	} else {
		if err := d.decode(name, data, reflect.Indirect(val)); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decodeFunc(name string, data interface{}, val reflect.Value) error {
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if val.Type() != dataVal.Type() {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}
	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeSlice(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	valSlice := val
	if valSlice.IsNil() {
		// Check input type
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {
			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		}

		// If the input value is empty, then don't allocate since non-nil != nil
		if dataVal.Len() == 0 {
			return nil
		}

		// Make a new slice to hold our result, same size as the original data.
		valSlice = reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	// Finally, set the value to the slice we built up
	val.Set(valSlice)

	// If there were errors, we return those
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeArray(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	if valArray.Interface() == reflect.Zero(valArray.Type()).Interface() {
		// Check input type
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {

			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		}
		if dataVal.Len() > arrayType.Len() {
			return fmt.Errorf(
				"'%s': expected source data to have length less or equal to %d, got %d",
				name, arrayType.Len(), dataVal.Len())

		}

		// Make a new array to hold our result, same size as the original data.
		valArray = reflect.New(arrayType).Elem()
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		currentField := valArray.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	// Finally, set the value to the array we built up
	val.Set(valArray)

	// If there were errors, we return those
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeStruct(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))

	// If the type of the value to write to and the data match directly,
	// then we just set it directly instead of recursing into the structure.
	if dataVal.Type() == val.Type() {
		val.Set(dataVal)
		return nil
	}

	switch dataVal.Kind() {
	case reflect.Map:
		return d.decodeStructFromMap(name, dataVal, val)
	case reflect.Struct:
		// Not the most efficient way to do this but we can optimize later if
		// we want to. To convert from struct to struct we go to map first
		// as an intermediary.
		m := make(map[string]interface{})
		mval := reflect.Indirect(reflect.ValueOf(&m))
		if err := d.decodeMapFromStruct(name, dataVal, mval, mval); err != nil {
			return err
		}

		return d.decodeStructFromMap(name, mval, val)
	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeStructFromMap(name string, dataVal, val reflect.Value) error {
	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	dataValKeys := make(map[reflect.Value]struct{})
	for _, dataValKey := range dataVal.MapKeys() {
		dataValKeys[dataValKey] = struct{}{}
	}

	errors := make([]string, 0)

	// Compile the list of all the fields that we're going to be decoding
	// from all the structs.
	type field struct {
		field reflect.StructField
		val   reflect.Value
	}
	var fields []field

	structType := val.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		// Normal struct field, store it away
		fields = append(fields, field{fieldType, val.Field(i)})
	}

	// for fieldType, field := range fields {
	for _, f := range fields {
		field, fieldValue := f.field, f.val
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" {
			fieldName = tagValue
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() {
			// Do a slower search by iterating over each key and
			// doing case-insensitive search.
			for dataValKey := range dataValKeys {
				mK, ok := dataValKey.Interface().(string)
				if !ok {
					// Not a string key
					continue
				}

				if strings.EqualFold(mK, fieldName) {
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				}
			}

			if !rawMapVal.IsValid() {
				// There was no matching key in the map for the value in
				// the struct. Just ignore.
				continue
			}
		}

		if !fieldValue.IsValid() {
			// This should never happen
			panic("field is not valid")
		}

		// If we can't set the field, then it is unexported or something,
		// and we just continue onwards.
		if !fieldValue.CanSet() {
			continue
		}

		// If the name is empty string, then we're at the root, and we
		// don't dot-join the fields.
		if name != "" {
			fieldName = fmt.Sprintf("%s.%s", name, fieldName)
		}

		if err := d.decode(fieldName, rawMapVal.Interface(), fieldValue); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch getKind(v) {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
