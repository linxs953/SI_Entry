package handler

import (
	"errors"
	"net/http"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getClusterCRsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetClusterCRsRequest
		// if err := httpx.Parse(r, &req); err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// 	return
		// }
		req.Type = r.URL.Query().Get("type")
		if req.Type == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("type is required"))
			return
		}
		l := logic.NewGetClusterCRsLogic(r.Context(), svcCtx)
		resp, err := l.GetClusterCRs(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
