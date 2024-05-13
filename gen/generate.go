package gen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cscoding21/csgen"
)

func Generate(file ...string) error {
	fullPath := csgen.GetFile()

	fmt.Println(fullPath)

	makeValidator := false

	structs, err := csgen.GetStructs("test_struct.go")
	if err != nil {
		return err
	}

	if len(structs) == 0 {
		return nil
	}

	pkg := structs[0].Package
	outFileName := "test_struct"

	builder := csgen.NewCSGenBuilderForFile("csval", pkg)
	builder.WriteString(getImportStatement())

	for _, st := range structs {
		fmt.Println(st.Name)
		fileContents := buildValidator(st, builder)

		for _, f := range st.Fields {
			valTags := f.GetTag("csval")

			if len(valTags) == 0 {
				continue
			}

			fmt.Print(fileContents)
			makeValidator = true
		}
	}

	if makeValidator {
		valFile := getGeneratedFileName(strings.ToLower(outFileName))
		err = csgen.WriteGeneratedGoFile(valFile, builder.String())
		if err != nil {
			fmt.Printf("error writing file: %v", err)
			return err
		}
	}

	return nil
}

func buildValidator(st csgen.Struct, builder *strings.Builder) string {
	builder.WriteString(fmt.Sprintf("func (obj *%s) Validate() csval.ValidationResult {", st.Name))
	builder.WriteByte('\n')

	builder.WriteString("result := csval.NewSuccessValidationResult()")
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
				builder.WriteString(getIsUrl(f.Name))
				builder.WriteByte('\n')
			}

			if _, ok := tm["ip"]; ok {
				builder.WriteString(getIsIP(f.Name))
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
		} else {
			if _, ok := tm["obj"]; ok {
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
			re := regexp.MustCompile(`^(\w*)\((\w*)\)$`)
			res := re.FindAllStringSubmatch(t, -1)
			key := res[0][1]
			val := res[0][2]

			tagMap[key] = val

		} else {
			tagMap[t] = 1
		}
	}

	return &tagMap
}

func getGeneratedFileName(originFile string) string {
	return fmt.Sprintf("%s_csval.gen.go", strings.TrimSuffix(originFile, ".go"))
}

func getImportStatement() string {
	return `import (
	csval "github.com/cscoding21/csval"
	)


	`
}

func getIsRequired(field string) string {
	return fmt.Sprintf("result.Append(csval.IsNotEmpty(\"%s\", obj.%s))", field, field)
}

func getIsEmail(field string) string {
	return fmt.Sprintf("result.Append(csval.IsEmail(\"%s\", obj.%s))", field, field)
}

func getIsUrl(field string) string {
	return fmt.Sprintf("result.Append(csval.IsValidWebAddress(\"%s\", obj.%s))", field, field)
}

func getIsIP(field string) string {
	return fmt.Sprintf("result.Append(csval.IsIP(\"%s\", obj.%s))", field, field)
}

func getIsGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(csval.IsGreaterThan(\"%s\", obj.%s, %v))", field, field, min)
}

func getIsLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(csval.IsLessThan(\"%s\", obj.%s, %v))", field, field, max)
}

func getIsLengthGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(csval.IsLengthGreaterThan(\"%s\", obj.%s, %v))", field, field, min)
}

func getIsLengthLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(csval.IsLengthLessThan(\"%s\", obj.%s, %v))", field, field, max)
}

func getCheckForObject(field string) string {
	return fmt.Sprintf("result.Append(obj.%s.Validate())", field)
}
