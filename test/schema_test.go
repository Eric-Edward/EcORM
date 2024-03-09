package schema

import (
	"EcORM/dialect"
	"fmt"
	"testing"
)

type User struct {
	Name string `pri:"primary key"`
	Age  int
}

var testDialect, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	result := Parse(&User{}, testDialect)

	if result.Name != "User" || len(result.filesMap) != 2 {
		t.Fatal("解析失败")
	}
	if result.GetField("Name").Tag != "primary key" {
		t.Fatal("标签解析失败")
	}
	fmt.Println(result)
	fmt.Println("解析成功")
}
