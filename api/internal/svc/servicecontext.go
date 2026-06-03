// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"api/internal/config"
	"rpc/user"

	//"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UserRpcClient user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	//logger.InitDefaultLogger(&c.LoggerConfig)
	//logx.SetWriter(logger.NewZapWriter(logger.GetZapLogger()))
	//logx.DisableStat()
	return &ServiceContext{
		Config:        c,
		UserRpcClient: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
