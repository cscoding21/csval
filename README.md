
<p align="center"><img src="https://github.com/cscoding21/cscoding/blob/main/assets/csc-banner.png?raw=true" width=728></p>

<p align="center">
    <a href="https://github.com/cscoding21/csval"><img src="https://img.shields.io/badge/built_with-Go-29BEB0.svg?style=flat-square"></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/cscoding21/csval"><img src="https://goreportcard.com/badge/github.com/cscoding21/csval?style=flat-square"></a>&nbsp;
 <a href="https://pkg.go.dev/mod/github.com/cscoding21/csval"><img src="https://pkg.go.dev/badge/mod/github.com/cscoding21/csval"></a>&nbsp;
    <a href="https://github.com/cscoding21/csval/" alt="Stars">
        <img src="https://img.shields.io/github/stars/cscoding21/csval?color=0052FF&labelColor=090422" /></a>&nbsp;
    <a href="https://github.com/cscoding21/csval/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/cscoding21/csval?color=0052FF&labelColor=090422" /></a>
    <br />
    <a href="https://discord.gg/BjV88Bys" alt="Discord">
        <img src="https://img.shields.io/discord/1196192809120710779" /></a>&nbsp;
    <a href="https://www.youtube.com/@CommonSenseCoding-ge5dn" alt="YouTube">
        <img src="https://img.shields.io/badge/youtube-watch_videos-red.svg?color=0052FF&labelColor=090422&logo=youtube" /></a>&nbsp;
    <a href="https://twitter.com/cscoding21" alt="YouTube">
        <img src="https://img.shields.io/twitter/follow/cscoding21" /></a>&nbsp;
</p>



# CSVal
CSVal is a Golang package that allows the developer to define object validation rules using struct tags, but without having to rely on runtime reflection, which can have suboptimal performance.  If reflection is not a problem...[Go Playground Validator](https://github.com/go-playground/validator) is an excellent alternative.  

CSVal consists of the following components:

- __generator__: a package that views target code and generated corresponding validation source.
- __validator__: a library of common validation rules that can be applied to the fields within structs.
- __runner__: an executable that runs via **go generate** and leverages the generator package against the developer's code

## Usage
To use CSVal, add the __csval__ tag with a comma-separated list of rules to your struct's field definitions.  For example...given the file __data.go__:

    //go:generate csval

    package tests

    // FooStruct is a sample struct to be validated
    type FooStruct struct {
        Name        string       `csval:"req"`
        Email       string       `csval:"req,email"`
        Password    string       `csval:"min(3),max(11)"`
        ConfirmPass string       `csval:"req,equals(Password)"`
        Age         int          `csval:"min(18),max(65)"`
        Sub         BarSubStruct `csval:"obj"`
    }

    // BarSubStruct is a sample struct to be validated
    type BarSubStruct struct {
        IP   string `csval:"req,ip"`
        Port int
    }

The above struct is annotated with csval tags that define validation rules for the object.  When run, the generator will create a new file in the same package called __data_csval.gen.go__.  The new file will contain a __Validate__ method with the __FooStruct__ struct as a receiver as follows.

    func (obj *FooStruct) Validate() validate.ValidationResult {
        result := validate.NewSuccessValidationResult()

        // ---Field: Name
        result.Append(validate.IsNotEmpty("Name", obj.Name))

        // ---Field: Email
        result.Append(validate.IsNotEmpty("Email", obj.Email))
        result.Append(validate.IsEmail("Email", obj.Email))

        // ---Field: Password
        result.Append(validate.IsLengthGreaterThan("Password", obj.Password, 3))
        result.Append(validate.IsLengthLessThan("Password", obj.Password, 11))

        // ---Field: ConfirmPass
        result.Append(validate.IsNotEmpty("ConfirmPass", obj.ConfirmPass))
        result.Append(validate.IsEqualTo("ConfirmPass:Password", obj.ConfirmPass, obj.Password))

        // ---Field: Age
        result.Append(validate.IsGreaterThan("Age", obj.Age, 18))
        result.Append(validate.IsLessThan("Age", obj.Age, 65))

        // ---Field: Sub
        result.Append(obj.Sub.Validate())

        return result
    }

The __ValidationResult__ object contains a pass/fail status and a list of error messages if applicable.

# Installation
To use CSVal within your Go project, there are two installation steps.  The first is to add the package:

    go get github.com/cscoding21/csval

Additionally, the runner needs to be installed on the target machine.  This executes the generator when __go generate__ is invoked.

    go install github.com/cscoding21/csval


## Usage
To use CSVal, annotate the struct fields with the __csval__ tag.  The following annotations are supported:
| Annotation | Description |
| --- | --- |
|__req__|The field is required|
|__email__|The field must be a valid email address|
|__ip__|The field must be a valid IP address|
|__url__|The field must be greater than or equal to x|
|__min(x)__|If the field type is string, the length must be greater than or equal to x.  If the type is numeric, the value must be greater than or equal to x|
|__max(x)__|If the field type is string, the length must be less than or equal to x.  If the type is numeric, the value must be less than or equal to x|
|__equals(x)__|The field must be equal to the value of another field, represented by x|
|__regex(x)__|The field must satisfy the regular expression represented by x|
|__validate__|The field will have its own Validate function called.  This is for fields of a type that has CS validation rules|

Post code generation, testing the validity of an object is as simple as calling the __Validate__ method.  For example...

    s := &FooStruct{}
    result := s.Validate()
    if !result.Pass {
        for _, err := range result.Messages {
            fmt.Println(err.Field, err.Message)
        }
    }

