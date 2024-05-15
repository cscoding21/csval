package csval

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
	Sub      TestSubStruct `csval:"obj"`
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
