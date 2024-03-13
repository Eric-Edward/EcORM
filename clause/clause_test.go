package clause

import (
	"fmt"
	"testing"
)

func TestClause(t *testing.T) {
	var clause Clause
	clause.Set(SELECT, []interface{}{"User2", []string{"*"}})
	clause.Set(WHERE, []interface{}{"Name = ?", "Eric"})
	clause.Set(ORDERBY, []interface{}{"Age Asc"})
	clause.Set(LIMIT, []interface{}{1})

	sql, vars := clause.Build([]Type{SELECT, WHERE, ORDERBY, LIMIT})
	fmt.Println(sql, vars)
}
