package session

import (
	"EcORM/clause"
	"EcORM/dialect"
	"EcORM/log"
	"EcORM/schema"
	"database/sql"
	"strings"
)

//在session这个函数中，我们可以根据函数中执行的sql语句的返回值来确定当前函数应该返回什么类型

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVal   []interface{}
	clause   clause.Clause
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) ClearSQL() {
	s.sql.Reset()
	s.sqlVal = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// Raw 函数，官方给的答案是...的形式，然后我们这里使用切片来进行表示，官方这样写就可以在甘肃中添加很多个参数，而我们就
// 需要将参数写成一个切片的形式
func (s *Session) Raw(sq string, value []interface{}) *Session {
	s.sql.WriteString(sq)
	s.sqlVal = append(s.sqlVal, value...)
	return s
}

// Exec 函数是用来实现用户的一些相关的操作的处理函数
func (s *Session) Exec() (sql.Result, error) {
	result, err := s.DB().Exec(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *Session) QueryRow() *sql.Row {
	query := s.DB().QueryRow(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	return query
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	query, err := s.DB().Query(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return query, nil
}
