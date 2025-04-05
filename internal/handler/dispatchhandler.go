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

type CommonResponse struct {
	Code    int
	Message string
	Data    interface{}
}

func dispatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DispatchResourceRequest
		response := types.DispatchResourceResponse{
			Code:    int(ie.Success),
			Message: "success",
		}

		if err := httpx.Parse(r, &req); err != nil {
			response.Code = http.StatusBadRequest
			response.Message = err.Error()
			httpx.OkJsonCtx(r.Context(), w, response)
			return
		}

		if req.Event == "createTask" && len(req.Spec) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			response.Code = http.StatusBadRequest
			response.Message = "参数校验失败: Need Specific Spec field"
			httpx.OkJsonCtx(r.Context(), w, response)
			return
		}

		if len(req.Metadata) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			response.Code = http.StatusBadRequest
			response.Message = "参数校验失败: Need Specific Metadata field"
			httpx.OkJsonCtx(r.Context(), w, response)
			return
		}

		l := logic.NewDispatchLogic(r.Context(), svcCtx)
		resp, err := l.Dispatch(&req)
		if err != nil {
			logx.Errorf("创建任务失败, err: %v", err)
			response.Code = int(ie.GetErrorCode(err))
			response.Message = err.Error()
		}
		response.Extra = resp.Extra
		httpx.OkJsonCtx(r.Context(), w, response)
	}
}
