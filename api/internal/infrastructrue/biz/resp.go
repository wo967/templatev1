package biz

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(ctx context.Context, w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			logx.Errorw("unknown error", logx.Field("error", err))
		}
		body.Code = int(e.Code())
		if e.Message() != "" {
			body.Msg += e.Message()
		}
		body.Data = struct{}{}
	} else {
		body.Msg = "OK"
		if resp == nil {
			resp = struct{}{}
		}
		body.Data = resp
	}
	httpx.OkJsonCtx(ctx, w, body)
}
