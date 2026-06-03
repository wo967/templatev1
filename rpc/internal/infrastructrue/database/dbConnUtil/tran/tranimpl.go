package tran

import (
	"rpc/internal/infrastructrue/database/dbConnUtil"
	"rpc/internal/infrastructrue/database/dbConnUtil/gorms"

	"gorm.io/gorm"
)

type TransactionImpl struct {
	conn dbConnUtil.DbConn
}

func (t *TransactionImpl) Action(f func(conn dbConnUtil.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction(db *gorm.DB) *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.New(db),
	}
}
