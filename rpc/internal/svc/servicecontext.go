package svc

import (
	"rpc/internal/config"
	"rpc/internal/infrastructrue/database"
	"rpc/internal/infrastructrue/database/dbConnUtil"
)

type ServiceContext struct {
	Config config.Config
	Db     *dbConnUtil.MsDB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Db:     database.ConnMysql(c.Mysql.DataSource),
	}
}
