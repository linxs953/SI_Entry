package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"
)

func dispatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DispatchResourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJsonCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDispatchLogic(r.Context(), svcCtx)
		resp, _ := l.Dispatch(&req)
		httpx.OkJsonCtx(r.Context(), w, resp)
		return
	}
}
