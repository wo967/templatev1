package dao

import (
	"context"
	"rpc/internal/infrastructrue/database/dao/model"
	"rpc/internal/infrastructrue/database/dbConnUtil/gorms"

	"gorm.io/gorm"
)

type UserDao struct {
	conn *gorms.GormConn
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		conn: gorms.New(db),
	}
}

func (m *UserDao) FindByUsername(ctx context.Context, username string) (user *model.User, err error) {
	session := m.conn.Session(ctx)
	u := new(model.User)
	err = session.Model(&model.User{}).
		Where("username = ?", username).Limit(1).
		Take(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}
