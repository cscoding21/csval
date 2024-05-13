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

	structs, err := csgen.GetStructs("test_struct.go")
	if err != nil {
		return err
	}

	for _, st := range structs {
		fmt.Println(st.Name)
		makeValidator := false

		for _, f := range st.Fields {
			valTags := f.GetTag("csval")

			if len(valTags) == 0 {
				continue
			}

			fileContents := buildValidator(st)

			fmt.Print(fileContents)

			valFile := getGeneratedFileName(strings.ToLower(st.Name))
			err = csgen.WriteGeneratedGoFile(valFile, fileContents)
			if err != nil {
				fmt.Printf("error writing file: %v", err)
			}
		}

		fmt.Printf("%s - create validator:\n\n %v", st.Name, makeValidator)
	}

	return nil
}

func buildValidator(st csgen.Struct) string {
	builder := csgen.NewCSGenBuilderForFile("csval", st.Package)

	builder.WriteString(getImportStatement())

	builder.WriteString("func (obj *TestStruct) Validate() csval.ValidationResult {")
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

		if _, ok := tm["req"]; ok {
			builder.WriteString(getIsRequired(f.Name))
			builder.WriteByte('\n')
		}

		if _, ok := tm["email"]; ok {
			builder.WriteString(getIsEmail(f.Name))
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
	return fmt.Sprintf("result.Append(csval.IsNotEmpty(obj.%s))", field)
}

func getIsEmail(field string) string {
	return fmt.Sprintf("result.Append(csval.IsEmail(obj.%s))", field)
}

func getIsGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(csval.IsGreaterThan(obj.%s, %v))", field, min)
}

func getIsLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(csval.IsLessThan(obj.%s, %v))", field, max)
}

func getIsLengthGreaterThan(field string, min string) string {
	return fmt.Sprintf("result.Append(csval.IsLengthGreaterThan(obj.%s, %v))", field, min)
}

func getIsLengthLessThan(field string, max string) string {
	return fmt.Sprintf("result.Append(csval.IsLengthLessThan(obj.%s, %v))", field, max)
}
