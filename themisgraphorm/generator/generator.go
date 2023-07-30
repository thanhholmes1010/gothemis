package generator

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type GraphEngineGeneratorImlp struct {
}

func (g *GraphEngineGeneratorImlp) Gen(rw *io.ReadWriter) {

}

func (g *GraphEngineGeneratorImlp) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = line[0 : len(line)-1]
		if err != nil {
			break
		}
		g.buildCommandWithQuery(line)
	}
}

func (g *GraphEngineGeneratorImlp) buildCommandWithQuery(line string) {
	str_arr := strings.Split(line, " ")
	if str_arr[0] != "gen" {
		return
	}
	if len(str_arr) > 3 {
		command := str_arr[1]
		packageName := strings.ToLower(str_arr[2])
		fieldArr := str_arr[3:]
		if command == "schema" {
			g.genWithSchema(packageName, fieldArr)
		}
	}
}

type FieldGenData struct {
	Name                     string
	GoType                   string
	SQLConstructorMethodName string
	SizeType                 string
}
type SchemaGenData struct {
	PackageName string
	SchemaName  string
	Fields      []*FieldGenData
}

var mapGenTypeToGoType = map[string]string{
	"Int":     "int32",
	"UInt":    "uint32",
	"Varchar": "string",
}

var mapGenTypeToSqlContructorMethodName = map[string]string{
	"Int":     "Integer",
	"UInt":    "Integer",
	"Varchar": "Varchar",
}

func (g *GraphEngineGeneratorImlp) genWithSchema(name string, fieldArr []string) {
	data := &SchemaGenData{
		PackageName: name,
		Fields:      make([]*FieldGenData, 0),
	}
	data.SchemaName = strings.ToUpper(string(name[0])) + name[1:]
	for _, field := range fieldArr {
		str_field := strings.Split(field, ":")
		fieldName := str_field[0]
		gen_type := ""
		size_type := ""
		get_size := false
		for _, char := range str_field[1] {
			if string(char) == "(" || string(char) == ")" {
				get_size = true
				continue
			}
			if !get_size {
				gen_type += string(char)
			} else {
				size_type += string(char)
			}
		}
		if gen_type == "Int" {
			size_type = "32"
		}
		fieldGenData := &FieldGenData{
			Name:                     fieldName,
			GoType:                   mapGenTypeToGoType[gen_type],
			SQLConstructorMethodName: mapGenTypeToSqlContructorMethodName[gen_type],
			SizeType:                 size_type,
		}
		fmt.Println(fieldGenData)
		data.Fields = append(data.Fields, fieldGenData)
	}
	if _, err := os.Stat(name); !os.IsExist(err) {
		os.Mkdir(name, 0777)
	}
	file_name := name + "/" + name + ".go"
	f, err := os.OpenFile(file_name, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Open or Create Failed on error: ", err)
	}
	if err == nil {
		// then generator
		tmpl := template.Must(template.New(name).Parse(string(genSchemaTemplate)))
		if err := tmpl.Execute(f, data); err != nil {
			fmt.Println("write error: ", err)
		}
	}

}

var genSchemaTemplate = []byte(`
package {{$.PackageName}}

type {{$.SchemaName}} struct {
{{- range $Field := $.Fields }}
    {{$Field.Name}} {{$Field.GoType -}}
{{ end }}
}

// DefineFields be generated from @themisgraphorm
func (sc *{{$.SchemaName}}) DefineFields() []field.IField {
     return []field.IField{
     {{- range $Field := $.Fields }}
         field.{{$Field.SQLConstructorMethodName }}({{$Field.SizeType -}}),
     {{ end }}
     }
}

// DefineEdges be generated from @themisgraphorm
func (sc *{{$.SchemaName}}) DefineEdges() []edge.IEdge {
     return []edge.IEdge{}
}
`)
