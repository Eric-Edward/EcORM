package session

import (
	"EcORM/dialect"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type User struct {
	Name string `pri:"primary key"`
	Age  int
}

var TestDialect, _ = dialect.GetDialect("mysql")

func TestCreateTable(t *testing.T) {
	db, _ := sql.Open("mysql", "root:Tsinghua@/EcORM")
	session := New(db, TestDialect)

	session.Model(&User{})
	_ = session.CreateTable()
	fmt.Println(session.sql.String())

}
