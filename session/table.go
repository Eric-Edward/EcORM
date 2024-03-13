package session

import (
	"EcORM/log"
	"EcORM/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
}

func (s *Session) GetRefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model的表还未设置")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var columns []string
	for _, field := range s.refTable.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	args := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s)", table.Name, args), nil).Exec()
	if err != nil {
		log.Error(fmt.Sprintf("创建表{%s}失败，原因如下：", table.Name), err)
	}
	return err
}

func (s *Session) DropTable() error {
	table := s.refTable
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name), nil).Exec()
	if err != nil {
		log.Error(fmt.Sprintf("表{%s}删除失败，原因如下：", table.Name), err)
	}
	s.refTable = nil
	return err
}

func (s *Session) TableExist() bool {
	name := s.refTable.Name
	_, err := s.Raw(s.dialect.TableExistSQL(name), nil).Exec()
	if err != nil {
		return false
	}
	return true
}
