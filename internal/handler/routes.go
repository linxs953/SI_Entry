// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"

	"entry/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/cluster/get",
				Handler: getClusterCRsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/dispatch",
				Handler: dispatchHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/from/:name",
				Handler: EntryHandler(serverCtx),
			},
		},
	)
}
