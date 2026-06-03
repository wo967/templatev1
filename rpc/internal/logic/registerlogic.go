package logic

import (
	"context"

	//"rpc/internal/infrastructrue/database/dao"
	"rpc/internal/svc"
	"rpc/template"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *RegisterLogic) Register(in *template.RegisterRequest) (*template.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	//user, err := dao.NewUserDao(l.svcCtx.Db.DB).FindByUsername(l.ctx, in.Username)
	return &template.RegisterResponse{}, nil
}
