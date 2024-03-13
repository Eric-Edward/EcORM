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

	fmt.Println(p)
	_ = session.DropTable()
}
