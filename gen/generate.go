package gen

import (
	"fmt"
	"os"
	"strings"

	"github.com/cscoding21/csgen"
)

// Generate creates validation methods for the structs in the file specified by go generate
func Generate(file ...string) error {
	fullPath := csgen.GetFile(file...)
	makeValidator := false

	structs, err := csgen.GetStructs(fullPath)
	if err != nil {
		return err
	}

	if len(structs) == 0 {
		return nil
	}

	pkg := structs[0].Package

	//---TODO: this is fragile and should be rethought
	outFileName := os.Getenv("GOFILE")

	builder := csgen.NewCSGenBuilderForFile("csval", pkg)

	for _, st := range structs {
		fmt.Println(st.Name)
		buildValidator(st, builder)

		for _, f := range st.Fields {
			valTags := f.GetTag("csval")

			if len(valTags) == 0 {
				continue
			}

			makeValidator = true
		}
	}

	if makeValidator {
		valFile := csgen.GetFileName("csval", "", outFileName)
		err = csgen.WriteGeneratedGoFile(valFile, builder.String())
		if err != nil {
			fmt.Printf("error writing file: %v", err)
			return err
		}
	}

	return nil
}

// CheckFileHasValidatorTags returns true if a file contains structs that have csvalidator tags and false otherwise
func CheckFileHasValidatorTags(path string) (bool, error) {
	structs, err := csgen.GetStructs(path)
	if err != nil {
		return false, err
	}

	for _, st := range structs {
		for _, field := range st.Fields {
			if len(field.GetTag("csval")) > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

func buildValidator(st csgen.Struct, builder *strings.Builder) string {
	builder.WriteString(fmt.Sprintf("// Validate checks the fields in the struct %s to ensure it conforms to business rules", st.Name))
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("func (obj *%s) Validate() validate.ValidationResult {", st.Name))
	builder.WriteByte('\n')

	builder.WriteString("result := validate.NewSuccessValidationResult()")
	builder.WriteByte('\n')
	builder.WriteByte('\n')

	for _, f := range st.Fields {
		valTags := f.GetTag("csval")
		tagMap := getTagMap(valTags)

		if tagMap == nil {
			continue
		}

		tm := *tagMap

		builder.WriteByte('\n')
		builder.WriteString(fmt.Sprintf("//---Field: %s", f.Name))
		builder.WriteByte('\n')

		if f.IsPrimitive {
			if _, ok := tm["req"]; ok {
				builder.WriteString(getIsRequired(f.Name))
				builder.WriteByte('\n')
			}

			if _, ok := tm["email"]; ok {
				builder.WriteString(getIsEmail(f.Name))
				builder.WriteByte('\n')
			}

			if _, ok := tm["url"]; ok {
				builder.WriteString(getIsURL(f.Name))
				builder.WriteByte('\n')
			}

			if _, ok := tm["ip"]; ok {
				builder.WriteString(getIsIP(f.Name))
				builder.WriteByte('\n')
			}

			if re, ok := tm["regex"]; ok {
				builder.WriteString(getSatifiesRegex(f.Name, re.(string)))
				builder.WriteByte('\n')
			}

			if min, ok := tm["min"]; ok {
				if f.Type == "string" {
					builder.WriteString(getIsLengthGreaterThan(f.Name, min.(string)))
				} else if f.Type == "int" {
					builder.WriteString(getIsGreaterThan(f.Name, min.(string)))
				}

				builder.WriteByte('\n')
			}

			if max, ok := tm["max"]; ok {
				if f.Type == "string" {
					builder.WriteString(getIsLengthLessThan(f.Name, max.(string)))
				} else if f.Type == "int" {
					builder.WriteString(getIsLessThan(f.Name, max.(string)))
				}

				builder.WriteByte('\n')
			}

			if equals, ok := tm["equals"]; ok {
				builder.WriteString(getIsEqualTo(f.Name, equals.(string)))
				builder.WriteByte('\n')
			}

		} else {
			if _, ok := tm["validate"]; ok {
				builder.WriteString(getCheckForObject(f.Name))
				builder.WriteByte('\n')
				continue
			}
		}
	}

	builder.WriteByte('\n')
	builder.WriteString("return result")
	builder.WriteByte('\n')

	builder.WriteString("}")
	builder.WriteByte('\n')
	builder.WriteByte('\n')

	return builder.String()
}

func getTagMap(tags string) *map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	tagMap := make(map[string]interface{})
	ta := strings.Split(tags, ",")

	for _, t := range ta {
		parens := strings.Contains(t, "(")

		if parens {
			fp := strings.Index(t, "(")
			key := t[0:fp]
			val := t[fp+1 : len(t)-1]

			tagMap[key] = val

		} else {
			tagMap[t] = 1
		}
	}

	return &tagMap
}

func getIsRequired(field string) string {
	return fmt.Sprintf("result.Append(validate.IsNotEmpty(\"%s\", obj.%s))", field, field)
}

func getIsEmail(field string) string {
	return fmt.Sprintf("result.Append(validate.IsEmail(\"%s\", obj.%s))", field, field)
}

func getIsURL(field string) string {
	return fmt.Sprintf("result.Append(validate.IsValidWebAddress(\"%s\", obj.%s))", field, field)
}

func getIsIP(field string) string {
	return fmt.Sprintf("result.Append(validate.IsIP(\"%s\", obj.%s))", field, field)
}

func getIsGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(validate.IsGreaterThan(\"%s\", obj.%s, %v))", field, field, min)
}

func getIsLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(validate.IsLessThan(\"%s\", obj.%s, %v))", field, field, max)
}

func getIsLengthGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(validate.IsLengthGreaterThan(\"%s\", obj.%s, %v))", field, field, min)
}

func getIsLengthLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(validate.IsLengthLessThan(\"%s\", obj.%s, %v))", field, field, max)
}

func getCheckForObject(field string) string {
	return fmt.Sprintf("result.Append(obj.%s.Validate())", field)
}

func getIsEqualTo(field1 string, field2 string) string {
	fi := fmt.Sprintf("%s:%s", field1, field2)
	return fmt.Sprintf("result.Append(validate.IsEqualTo(\"%s\", obj.%s, obj.%s))", fi, field1, field2)
}

func getSatifiesRegex(field string, re string) string {
	return fmt.Sprintf("result.Append(validate.SatisfiesRegex(\"%s\", obj.%s, \"%s\"))", field, field, re)
}
