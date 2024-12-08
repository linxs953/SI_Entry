package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"entry/internal/proto/scheduler"
	"entry/internal/svc"
	"entry/internal/types"
)

type FetchResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchResourceLogic {
	return &FetchResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchResourceLogic) getTaskDefine(client scheduler.SchedulerClient, namespace, name string) (map[string]interface{}, error) {
	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		return nil, errors.New("name is empty")
	}
	grpcResp, err := client.GetTaskDefine(l.ctx, &scheduler.GetTaskDefineRequest{
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("call scheduler service failed: %v", err)
	}

	if grpcResp.Data == nil {
		return nil, fmt.Errorf("TaskDefine obj is nil")
	}

	data := map[string]interface{}{
		"spec":   grpcResp.Data.Spec,
		"status": grpcResp.Data.Status,
	}
	return data, nil
}

func (l *FetchResourceLogic) FetchResource(req *types.FetchResourceRequest) (resp *types.FetchResourceResponse, err error) {
	resp = &types.FetchResourceResponse{}

	// 获取gRPC连接
	conn, ok := l.svcCtx.ClientManager.GetClient("scheduler")
	if !ok {
		return nil, errors.New("scheduler service not found")
	}

	// 创建scheduler client
	client := scheduler.NewSchedulerClient(conn)

	switch req.ResourceType {
	case "1": // TaskDefine type
		data, err := l.getTaskDefine(client, req.Namespace, req.ResourceName)
		if err != nil {
			resp.Code = 1
			resp.Message = err.Error()
			return resp, err
			// return nil, err
		}
		resp.Data = data
		resp.Code = 0
		resp.Message = "success"
		return resp, err
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", req.ResourceType)
	}

}
