package model

import "time"

// 用户表
type User struct {
	Id            int64     `gorm:"column:id"`
	Username      string    `gorm:"column:username"`
	Password      string    `gorm:"column:password"`
	RegisterTime  time.Time `gorm:"column:register_time"`
	LastLoginTime time.Time `gorm:"column:last_login_time"`
}

func (*User) TableName() string {
	return "user"
}

// 做一些处理，例如前端传入的与mysql中的类型不一致，做一些转换
func (m *User) Superxxx(partner string) {
}
