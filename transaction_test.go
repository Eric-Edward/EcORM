package EcORM

import (
	"EcORM/dialect"
	"EcORM/log"
	"EcORM/session"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var testDialect, _ = dialect.GetDialect("mysql")

type User struct {
	Name string
	Age  int
}

type Student struct {
	Name   string
	Number int
}

func TestTransaction(t *testing.T) {
	e, _ := NewEngine("mysql", "root:Tsinghua@/EcORM")
	defer func() { _ = e.db.Close() }()
	_, err := e.Transaction(func(s *session.Session) ([]interface{}, error) {
		s.Model(&User{})
		err := s.CreateTable()
		if err != nil {
			return nil, err
		}
		_, err2 := s.Insert(User{Name: "Eric", Age: 23}, User{Name: "Ray", Age: 20})
		if err2 != nil {
			return nil, err2
		}
		//_ = s.DropTable()
		return nil, err
	})
	if err != nil {
		return
	}
}

func TestEngine_Migrate(t *testing.T) {
	e, _ := NewEngine("mysql", "root:Tsinghua@/EcORM")
	e.Transaction(func(s *session.Session) ([]interface{}, error) {
		s.Model(&User{})
		err1 := s.CreateTable()
		u1 := User{Name: "Eric", Age: 23}
		u2 := User{Name: "wda", Age: 22}
		_, err2 := s.Insert(u1, u2)
		return nil, errors.Join(err1, err2)
	})
	err := e.Migrate(User{}, Student{})
	if err != nil {
		log.Error("迁移出现问题，原因如下:", err)
	}
}
