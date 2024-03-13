package session

import (
	"EcORM/clause"
	"EcORM/dialect"
	"EcORM/log"
	"EcORM/schema"
	"database/sql"
	"reflect"
	"strings"
)

//在session这个函数中，我们可以根据函数中执行的sql语句的返回值来确定当前函数应该返回什么类型

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVal   []interface{}
	clause   clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) ClearSQL() {
	s.sql.Reset()
	s.sqlVal = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() *sql.DB {
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
	result, err := s.db.Exec(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *Session) QueryRow() *sql.Row {
	query := s.db.QueryRow(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	return query
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	query, err := s.db.Query(s.sql.String(), s.sqlVal...)
	defer s.ClearSQL()
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return query, nil
}

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.Model(value)
		table := s.GetRefTable()
		if !s.clause.IsSet() {
			s.clause.Set(clause.INSERT, []interface{}{table.Name, table.FieldsName})
		}
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues)
	sq, vars := s.clause.Build([]clause.Type{clause.INSERT, clause.VALUES})
	result, err := s.Raw(sq, vars).Exec()
	if err != nil {
		log.Error("插入执行出错，原因如下：", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(value interface{}) error {
	descSlice := reflect.Indirect(reflect.ValueOf(value))
	descType := descSlice.Type().Elem()
	s.Model(descType)
	s.clause.Set(clause.SELECT, []interface{}{s.refTable.Name, s.refTable.FieldsName})
	sq, vars := s.clause.Build([]clause.Type{clause.SELECT})
	result, err := s.Raw(sq, vars).QueryRows()
	if err != nil {
		log.Error("查询失败，原因如下:", err)
		return err
	}

	for result.Next() {
		descNumber := reflect.New(descType).Elem()
		var values []interface{}
		for _, field := range s.refTable.FieldsName {
			values = append(values, descNumber.FieldByName(field).Addr().Interface())
		}
		if err := result.Scan(values...); err != nil {
			log.Error("查询结果匹配对象失败，具体原因如下:", err)
			return err
		}
		descSlice.Set(reflect.Append(descSlice, descNumber))
	}
	return result.Close()
}
