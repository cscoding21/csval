package validate

import (
	"fmt"
	"testing"
)

// ---Sample of struct to be validated
type TestStruct struct {
	Name     string `csval:"req"`
	Email    string `csval:"req,email"`
	Address  string
	Password string        `csval:"min(3),max(11)"`
	Age      int           `csval:"min(18),max(65)"`
	Sub      TestSubStruct `csval:"validate"`
}

type TestSubStruct struct {
	IP   string `csval:"req,ip"`
	Port int
}

// ---Sample test data
var testData = TestStruct{
	Name:     "test",
	Email:    "tet@Test.com",
	Address:  "",
	Password: "password",
	Age:      21,
	Sub: TestSubStruct{
		IP:   "0.0.0.1",
		Port: 8080,
	},
}

// ---#########################################################################
// ---Sample of validate function created by code generator
func (obj *TestStruct) Validate() ValidationResult {
	result := NewSuccessValidationResult()

	// ---Field: Name
	result.Append(IsNotEmpty("Name", obj.Name))

	// ---Field: Email
	result.Append(IsNotEmpty("Email", obj.Email))
	result.Append(IsEmail("Email", obj.Email))

	// ---Field: Password
	result.Append(IsLengthGreaterThan("Password", obj.Password, 3))
	result.Append(IsLengthLessThan("Password", obj.Password, 11))

	// ---Field: Age
	result.Append(IsGreaterThan("Age", obj.Age, 18))
	result.Append(IsLessThan("Age", obj.Age, 65))

	// ---Field: Sub
	result.Append(obj.Sub.Validate())

	return result
}

func (obj *TestSubStruct) Validate() ValidationResult {
	result := NewSuccessValidationResult()

	// ---Field: IP
	result.Append(IsNotEmpty("IP", obj.IP))
	result.Append(IsIP("IP", obj.IP))

	return result
}

//---#########################################################################

func TestObjectValidate(t *testing.T) {
	result := testData.Validate()

	t.Log(result)

	if result.Pass != true {
		t.Errorf("object validation failed")
	}
}

func TestGetAggregatedError(t *testing.T) {
	testStruct := TestStruct{
		Name:  "test",
		Email: "tet@Test.com",
	}

	result := testStruct.Validate()

	if result.Pass {
		t.Error("validation should fail")
	}

	err := result.Error("TestGetAggregatedError unit test function failed")
	if err == nil {
		t.Error("error aggregator returned nil value. expected messages")
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

	for i, input := range testCases {
		if IsNotEmpty(fmt.Sprintf("test%v", i), input.have).Pass != input.want {
			t.Errorf("string Required failed with input %s", input.have)
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

	for i, input := range testCases {
		result := IsEmail(fmt.Sprintf("test_email_%v", i), input.have)
		if result.Pass != input.want {
			t.Errorf("Email failed with input %s", input.have)
		}
	}
}

func TestSatisfyRegex(t *testing.T) {
	testCases := []struct {
		ok    bool
		have  string
		regex string
		want  bool
	}{
		{ok: true, have: "abc21", regex: "^[a-z0-9_-]{3,16}$", want: true},
		{ok: true, have: "x2", regex: "^[a-z0-9_-]{3,16}$", want: false},
		{ok: true, have: "abc21&", regex: "^[a-z0-9_-]{3,16}$", want: false},
		{ok: true, have: "a", regex: "[a-z]", want: true},
		{ok: true, have: "X", regex: "[a-z]", want: false},
		{ok: true, have: "9", regex: "[a-z]", want: false},
		{ok: true, have: "blah", regex: "^(0?[1-9]|1[0-2]):[0-5][0-9]$", want: false},
		{ok: true, have: "09:59", regex: "^(0?[1-9]|1[0-2]):[0-5][0-9]$", want: true},
	}

	for _, input := range testCases {
		result := SatisfiesRegex(input.have, input.regex)
		if result.Pass != input.want {
			t.Errorf("Regex %s failed with input %s", input.regex, input.have)
		}
	}
}
