package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
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

		req.Event = r.URL.Query().Get("event")

		l := logic.NewFetchResourcesLogic(r.Context(), svcCtx)
		resp, err := l.FetchResources(&req)
		if err != nil {
			logx.Errorf("获取资源列表失败, err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		}
		if resp.Code > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
