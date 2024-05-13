package csval

import (
	"testing"
)

// ---Sample of struct to be validated
type TestStruct struct {
	Name    string `csval:"req"`
	Email   string `csval:"req,email"`
	Address string
	Age     int `csval:"min(18),max(65)"`
}

// ---Sample test data
var testData = TestStruct{
	Name:    "test",
	Email:   "tet@Test.com",
	Address: "",
	Age:     21,
}

// ---Sample of validate function created by code generator
func (obj *TestStruct) Validate() ValidationResult {
	result := NewSuccessValidationResult()

	//---Name
	result.Append(IsNotEmpty(obj.Name))

	//---Email
	result.Append(IsNotEmpty(obj.Email))
	result.Append(IsEmail(obj.Email))

	//---Age
	result.Append(IsGreaterThan(obj.Age, 18))
	result.Append(IsLessThan(obj.Age, 65))

	return result
}

func TestObjectValidate(t *testing.T) {
	result := testData.Validate()

	t.Log(result)

	if result.Pass != true {
		t.Errorf("object validation failed")
	}
}

func TestIsNotEmpty(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want bool
	}{
		{ok: true, have: "abc", want: true},
		{ok: true, have: "", want: false},
	}

	for _, input := range testCases {
		if IsNotEmpty(input.have).Pass != input.want {
			t.Errorf("stringRequired failed with input %s", input.have)
		}
	}
}

func TestIsEmail(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want bool
	}{
		{ok: true, have: "abc", want: false},
		{ok: true, have: "jeph@cscoding.io", want: true},
	}

	for _, input := range testCases {
		if IsEmail(input.have).Pass != input.want {
			t.Errorf("stringRequired failed with input %s", input.have)
		}
	}
}
