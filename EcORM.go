package EcORM

import (
	"EcORM/dialect"
	"EcORM/log"
	"EcORM/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

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
