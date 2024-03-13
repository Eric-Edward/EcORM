package dialect

import (
	"reflect"
)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) string
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectsMap[name]
	return dialect, ok
}
