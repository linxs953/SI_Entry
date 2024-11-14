package handler

import (
	"net/http"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func EntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewEntryLogic(r.Context(), svcCtx)
		resp, err := l.Entry(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
