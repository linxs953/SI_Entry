// Code generated by goctl. DO NOT EDIT.
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
				Path:    "/resource/fetch",
				Handler: fetchResourceHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/resource/fetchAll",
				Handler: fetchResourcesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/resource/dispatch",
				Handler: dispatchHandler(serverCtx),
			},
		},
	)
}
