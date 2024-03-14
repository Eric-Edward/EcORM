package session

import (
	"EcORM/dialect"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type User struct {
	Name string
	Age  int
}

var testDialect, _ = dialect.GetDialect("mysql")

func TestSession_Insert(t *testing.T) {
	db, _ := sql.Open("mysql", "root:Tsinghua@/EcORM")
	session := New(db, testDialect)
	session.Model(User{})
	_ = session.CreateTable()

	u1 := User{Name: "eric", Age: 23}
	u2 := User{Name: "wda", Age: 23}

	_, _ = session.Insert(u1, u2)
	var p []User

	_ = session.Find(&p)

	upd := make(map[string]interface{})
	upd["Name"] = "hhhhh"
	_, _ = session.Where("Name='Eric'").Update(upd)
	_, _ = session.Where("Name='wda'").Delete()
	count, _ := session.Count()

	fmt.Println(p, count)
	_ = session.DropTable()
}
