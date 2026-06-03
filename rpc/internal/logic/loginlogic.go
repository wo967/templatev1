package logic

import (
	"context"

	"rpc/internal/infrastructrue/database/dao"
	"rpc/internal/svc"
	"rpc/template"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *LoginLogic) Login(in *template.LoginRequest) (*template.LoginResponse, error) {
	// todo: add your logic here and delete this line
	userDao := dao.NewUserDao(l.svcCtx.Db.DB)
	user, err := userDao.FindByUsername(l.ctx, in.Username)
	l.Logger.Info("用户登录", user)

	if err != nil {
		return nil, err
	}
	return &template.LoginResponse{
		Token: "123",
	}, nil
}
