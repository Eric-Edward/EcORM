package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	WHERE
	ORDERBY
	LIMIT
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(t Type, values []interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	c.sql[t], c.sqlVars[t] = generators[t](values)
}

func (c *Clause) Build(orders []Type) (string, []interface{}) {
	var sq []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sq = append(sq, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sq, " "), vars
}

func (c *Clause) IsSet() bool {
	return !(c.sql == nil)
}
