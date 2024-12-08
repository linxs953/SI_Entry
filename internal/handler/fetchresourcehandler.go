package handler

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"
)

func fetchResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FetchResourceRequest
		// if err := httpx.Parse(r, &req); err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// 	return
		// }

		req.ResourceName = r.URL.Query().Get("name")
		req.ResourceType = r.URL.Query().Get("type")
		if req.ResourceName == "" || req.ResourceType == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("name or type is empty"))
			return
		}

		l := logic.NewFetchResourceLogic(r.Context(), svcCtx)
		resp, _ := l.FetchResource(&req)
		httpx.OkJsonCtx(r.Context(), w, resp)

	}
}
