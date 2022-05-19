package schema

import (
	"geego/dialect"
	"reflect"
)

type Filed struct {
	Name string

	Type string

	Tag string
}

type Schema struct {
	Model      interface{}       //被映射的对象 Model
	Name       string            //表名 Name
	Fileds     []*Filed          //字段 Fields
	FiledNames []string          //包含所有的字段名(列名)
	filedMap   map[string]*Filed //记录字段名和 Field 的映射关系
}

func (schema *Schema) GetField(name string) *Filed {
	return schema.filedMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		filedMap: make(map[string]*Filed),
	}

}
