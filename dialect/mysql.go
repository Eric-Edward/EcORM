package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

func (m mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Int8, reflect.Uint8:
		return "TINYINT"
	case reflect.Int16, reflect.Uint16:
		return "SMALLINT"
	case reflect.Int32, reflect.Uint32:
		return "INTEGER"
	case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "TEXT"
	case reflect.Bool:
		return "TINYBLOB"
	case reflect.Array, reflect.Slice:
		return "BLOB"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "DATETIME"
		}
	default:
		panic(fmt.Sprintf("Invalid type:%s (%s) in mysql", typ.Type().Name(), typ.Kind()))
	}
	return ""
}

func (m mysql) TableExistSQL(tableName string) string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
}

// 这个是用来检测mysql有没有完全实现Dialect的方法
var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}
