package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	ie "entry/internal/error"
	"entry/internal/logic"
	"entry/internal/svc"
	"entry/internal/types"
)

func fetchResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FetchResourceRequest
		req.ResourceName = r.URL.Query().Get("name")
		req.ResourceType = r.URL.Query().Get("type")
		req.Event = r.URL.Query().Get("event")
		req.Metadata = make(map[string]interface{})
		taskId := r.URL.Query().Get("taskId")
		if taskId == "" {
			httpx.OkJsonCtx(r.Context(), w, &types.FetchResourceResponse{
				Code:    int(ie.InvalidParameter),
				Message: "name or type or taskId is empty",
			})
			return
		}
		req.Metadata["taskId"] = taskId

		l := logic.NewFetchResourceLogic(r.Context(), svcCtx)
		resp, err := l.FetchResource(&req)
		if err != nil {
			logx.Errorf("获取资源失败, err: %v", err)
			// httpx.ErrorCtx(r.Context(), w, ie.NewError(ie.InternalError, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		}
		if resp.Code == 100 {
			resp.Code = int(ie.Success)
		}
		if resp.Code > 100 {
			w.WriteHeader(http.StatusInternalServerError)
			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
