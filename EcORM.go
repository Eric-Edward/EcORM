package EcORM

import (
	"EcORM/dialect"
	"EcORM/log"
	"EcORM/session"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

type TxFunc func(s *session.Session) ([]interface{}, error)

func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	log.Info(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	d, _ := dialect.GetDialect(driver)

	//如果ping失败了，还是证明数据没有连接成功
	if e := db.Ping(); e != nil {
		log.Error(err)
		return nil, e
	}

	log.Info("连接数据库成功")

	return &Engine{
		db:      db,
		dialect: d,
	}, nil
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error("关闭数据库失败")
		return
	}
	log.Info("关闭数据库成功")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

func (e *Engine) Transaction(f TxFunc) (result []interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		log.Error("事务开启失败:", err)
		return nil, err
	}
	_, _ = s.DB().Exec("SET AUTOCOMMIT=0")
	defer func() {
		if p := recover(); p != nil {
			_ = s.RollBack()
			panic(p)
		} else if err != nil {
			_ = s.RollBack()
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}

func diff(a, b []string) (diff []string) {
	m := make(map[string]bool)
	for _, v := range b {
		m[v] = true
	}
	for _, v := range a {
		if _, ok := m[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

func (e *Engine) Migrate(value ...interface{}) error {
	_, err := e.Transaction(func(s *session.Session) ([]interface{}, error) {
		tableName := reflect.TypeOf(value[0]).Name()
		s.Model(value[1])
		if exist := s.TableExist(); !exist {
			_ = s.CreateTable()
		}

		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName), nil).QueryRows()
		columns, _ := rows.Columns()
		_ = rows.Close()

		addColumus := diff(s.GetRefTable().FieldsName, columns)
		delColumus := diff(columns, s.GetRefTable().FieldsName)

		for _, columu := range addColumus {
			field := s.GetRefTable().GetField(columu)
			_, err := s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, field.Name, field.Type), nil).Exec()
			if err != nil {
				return nil, errors.New(err.Error())
			}
		}

		if len(delColumus) == 0 {
			return nil, nil
		}
		tmp := "tmp_" + tableName
		fields := strings.Join(s.GetRefTable().FieldsName, ",")
		_, err1 := s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s", tmp, fields, tableName), nil).Exec()
		_, err2 := s.Raw(fmt.Sprintf("DROP TABLE %s", tableName), nil).Exec()
		_, err3 := s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME AS %s", tmp, s.GetRefTable().Name), nil).Exec()
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, errors.Join(err1, err2, err3)
		}
		return nil, nil
	})
	return err
}
