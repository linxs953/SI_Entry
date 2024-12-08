package handler

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"
)

func fetchResourcesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FetchAllResourcesRequest
		req.ResourceType = r.URL.Query().Get("type")
		req.Namespace = r.URL.Query().Get("namespace")
		if req.Namespace == "" {
			req.Namespace = "default"
		}
		if req.ResourceType == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("type is empty"))
			return
		}

		l := logic.NewFetchResourcesLogic(r.Context(), svcCtx)
		resp, err := l.FetchResources(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
