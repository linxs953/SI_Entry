package handler

import (
	"net/http"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func dispatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DispatchRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDispatchLogic(r.Context(), svcCtx)
		resp, err := l.Dispatch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
