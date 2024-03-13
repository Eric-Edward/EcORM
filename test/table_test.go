package main

import (
	"EcORM/dialect"
	"EcORM/session"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type User2 struct {
	Name string
	Age  int
}

var TestDialect, _ = dialect.GetDialect("mysql")

func TestCreateTable(t *testing.T) {
	db, _ := sql.Open("mysql", "root:Tsinghua@/EcORM")
	s := session.New(db, TestDialect)

	s.Model(&User2{})
	_ = s.CreateTable()
	fmt.Println(s.TableExist())
}
