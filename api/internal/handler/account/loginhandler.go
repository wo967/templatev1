// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"net/http"

	"api/internal/infrastructrue/biz"
	"api/internal/logic/account"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			biz.Response(r.Context(), w, nil, err)
			return
		}

		l := account.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		biz.Response(r.Context(), w, resp, err)
	}
}
