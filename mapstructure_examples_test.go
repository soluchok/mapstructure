package mapstructure

import (
	"fmt"
)

func ExampleDecode() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)
	// Output:
	// mapstructure.Person{Name:"Mitchell", Age:91, Emails:[]string{"one", "two", "three"}, Extra:map[string]string{"twitter":"mitchellh"}}
}

func ExampleDecode_errors() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   123,
		"age":    "bad value",
		"emails": []int{1, 2, 3},
	}

	var result Person
	err := Decode(input, &result)
	if err == nil {
		panic("should have an error")
	}

	fmt.Println(err.Error())
	// Output:
	// 5 error(s) decoding:
	//
	// * 'Age' expected type 'int', got unconvertible type 'string'
	// * 'Emails[0]' expected type 'string', got unconvertible type 'int'
	// * 'Emails[1]' expected type 'string', got unconvertible type 'int'
	// * 'Emails[2]' expected type 'string', got unconvertible type 'int'
	// * 'Name' expected type 'string', got unconvertible type 'int'
}

func ExampleDecode_tags() {
	// Note that the mapstructure tags defined in the struct type
	// can indicate which fields the values are mapped to.
	type Person struct {
		Name string `json:"person_name"`
		Age  int    `json:"person_age"`
	}

	input := map[string]interface{}{
		"person_name": "Mitchell",
		"person_age":  91,
	}

	var result Person
	err := Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)
	// Output:
	// mapstructure.Person{Name:"Mitchell", Age:91}
}

func ExampleDecode_omitempty() {
	// Add omitempty annotation to avoid map keys for empty values
	type Family struct {
		LastName string
	}
	type Location struct {
		City string
	}
	type Person struct {
		*Family   `json:",omitempty"`
		*Location `json:",omitempty"`
		Age       int
		FirstName string
	}

	result := &map[string]interface{}{}
	input := Person{FirstName: "Somebody"}
	err := Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result)
	// Output:
	// &map[Age:0 FirstName:Somebody]
}
