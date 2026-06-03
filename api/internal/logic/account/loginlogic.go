// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"
	"rpc/template"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	// 基础校验

	// 调用rpc
	rpcResp, err := l.svcCtx.UserRpcClient.Login(l.ctx, &template.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	//l.Logger.Info("rpc error")
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: rpcResp.Token,
	}, err
}
