package session

import "EcORM/log"

func (s *Session) Begin() (err error) {
	log.Info("事务开始")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error("事务开启失败:", err)
		return err
	}
	return nil
}

func (s *Session) Commit() error {
	log.Info("事务提交处理")
	if err := s.tx.Commit(); err != nil {
		log.Error("事务提交失败:", err)
		return err
	}
	return nil
}

func (s *Session) RollBack() error {
	log.Info("事务回滚")
	if err := s.tx.Rollback(); err != nil {
		log.Error("事务回滚失败:", err)
		return err
	}
	return nil
}
