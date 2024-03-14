package clause

import (
	"fmt"
	"strings"
)

type generator func(values []interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[VALUES] = _values
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[LIMIT] = _limit
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindStr(num int) string {
	var sql []string
	for i := 0; i < num; i++ {
		sql = append(sql, "?")
	}
	return strings.Join(sql, ",")
}

func _insert(values []interface{}) (string, []interface{}) {
	tableName := values[0]
	field := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %v (%v)", tableName, field), nil
}

func _select(values []interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), nil
}

func _values(values []interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	var bindStr string

	sql.WriteString("VALUES")

	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindStr(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(",")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _where(values []interface{}) (string, []interface{}) {
	return fmt.Sprintf("WHERE %s", values[0]), nil
}

func _orderBy(values []interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), nil
}

func _limit(values []interface{}) (string, []interface{}) {
	return fmt.Sprintf("LIMIT %v", values[0]), nil
}

func _update(values []interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+"=?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ",")), vars
}

func _delete(values []interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), nil
}

func _count(values []interface{}) (string, []interface{}) {
	vars := []interface{}{
		values[0],
		[]string{"count(*)"},
	}
	return _select(vars)
}
