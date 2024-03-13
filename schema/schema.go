package schema

import (
	"EcORM/dialect"
	"fmt"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{} //结构体的类型
	Name       string      //结构体的名字
	Fields     []*Field    //结构体的字段
	FieldsName []string
	filesMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.filesMap[name]
}

func (s *Schema) String() string {
	return fmt.Sprintf("Model:%s\nName:%s\nFilesName:%s", s.Model, s.Name, s.FieldsName)
}

func Parse(src interface{}, dialect dialect.Dialect) *Schema {
	model := reflect.Indirect(reflect.ValueOf(src)).Type()
	schema := &Schema{
		Model:    model,
		Name:     model.Name(),
		filesMap: make(map[string]*Field),
	}

	for i := 0; i < model.NumField(); i++ {
		f := model.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			file := &Field{
				Name: f.Name,
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			if value, ok := f.Tag.Lookup("pri"); ok {
				file.Tag = value
			}
			schema.Fields = append(schema.Fields, file)
			schema.FieldsName = append(schema.FieldsName, file.Name)
			schema.filesMap[file.Name] = file
		}
	}

	return schema
}

func (s *Schema) RecordValues(values interface{}) []interface{} {
	value := reflect.Indirect(reflect.ValueOf(values))
	var fieldValue []interface{}
	for _, field := range s.Fields {
		fieldValue = append(fieldValue, value.FieldByName(field.Name).Interface())
	}
	return fieldValue
}
