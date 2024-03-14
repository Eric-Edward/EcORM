package session

import (
	"EcORM/clause"
	"EcORM/log"
	"reflect"
)

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
	s.clause.SetType(clause.INSERT)
	s.clause.SetType(clause.VALUES)

	usedType := make([]clause.Type, 0)

	for _, typ := range *s.clause.GetUsedType() {
		usedType = append(usedType, typ)
	}
	sq, v := s.clause.Build(usedType)
	result, err := s.Raw(sq, v).Exec()
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
	s.clause.SetType(clause.SELECT)

	usedType := make([]clause.Type, 0)
	typ := *s.clause.GetUsedType()
	for i := len(typ) - 1; i >= 0; i-- {
		usedType = append(usedType, typ[i])
	}
	sq, v := s.clause.Build(usedType)
	result, err := s.Raw(sq, v).QueryRows()
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

func (s *Session) Update(values ...interface{}) (int64, error) {
	kv := values[0].(map[string]interface{})
	s.clause.Set(clause.UPDATE, []interface{}{s.refTable.Name, kv})
	s.clause.SetType(clause.UPDATE)

	usedType := make([]clause.Type, 0)
	typ := *s.clause.GetUsedType()
	for i := len(typ) - 1; i >= 0; i-- {
		usedType = append(usedType, typ[i])
	}
	sq, v := s.clause.Build(usedType)

	result, err := s.Raw(sq, v).Exec()
	if err != nil {
		log.Error("更新信息时产生错误，原因如下:", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, []interface{}{s.refTable.Name})
	s.clause.SetType(clause.DELETE)
	usedType := make([]clause.Type, 0)
	typ := *s.clause.GetUsedType()
	for i := len(typ) - 1; i >= 0; i-- {
		usedType = append(usedType, typ[i])
	}
	sq, v := s.clause.Build(usedType)
	result, err := s.Raw(sq, v).Exec()
	if err != nil {
		log.Error("删除信息出错，原因如下:", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, []interface{}{s.refTable.Name})
	s.clause.SetType(clause.COUNT)
	usedType := make([]clause.Type, 0)
	typ := *s.clause.GetUsedType()
	for i := len(typ) - 1; i >= 0; i-- {
		usedType = append(usedType, typ[i])
	}
	sq, v := s.clause.Build(usedType)
	result, err := s.Raw(sq, v).Exec()
	if err != nil {
		log.Error("计数发生错误，原因如下:", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Where(value interface{}) *Session {
	s.clause.Set(clause.WHERE, []interface{}{value.(string)})
	s.clause.SetType(clause.WHERE)
	return s
}

func (s *Session) ORDERBY(value interface{}) *Session {
	s.clause.Set(clause.ORDERBY, []interface{}{value})
	s.clause.SetType(clause.ORDERBY)
	return s
}

func (s *Session) Limit(value interface{}) *Session {
	s.clause.Set(clause.LIMIT, []interface{}{value.(int)})
	s.clause.SetType(clause.LIMIT)
	return s
}
