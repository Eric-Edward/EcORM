package session

import (
	"database/sql"
	"strings"
)

type Session struct {
	db     *sql.DB
	sql    strings.Builder
	sqlVal []interface{}
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}

func ClearSQL(s *Session) {
	s.sql.Reset()
	s.sqlVal = nil
}

func DB(s *Session) *sql.DB {
	return s.db
}

// Raw 函数，官方给的答案是...的形式，然后我们这里使用切片来进行表示，官方这样写就可以在甘肃中添加很多个参数，而我们就
// 需要将参数写成一个切片的形式
func (s *Session) Raw(sq string, value []interface{}) *Session {
	s.sql.WriteString(sq)
	s.sqlVal = append(s.sqlVal, value)
	return s
}
