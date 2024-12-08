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

type FetchResourcesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchResourcesLogic {
	return &FetchResourcesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchResourcesLogic) FetchResources(req *types.FetchAllResourcesRequest) (resp *types.FetchAllResourcesResponse, err error) {
	// todo: add your logic here and delete this line
	resp = &types.FetchAllResourcesResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	}
	// 获取gRPC连接
	conn, ok := l.svcCtx.ClientManager.GetClient("scheduler")
	if !ok {
		return nil, errors.New("scheduler service not found")
	}

	// 创建scheduler client
	client := scheduler.NewSchedulerClient(conn)
	switch req.ResourceType {
	case "1": // TaskDefine type
		data, err := l.getAllTaskDefine(client, req.Namespace)
		if err != nil {
			resp.Code = 1
			resp.Message = err.Error()
			return resp, err
		}
		resp.Data = data
		resp.Code = 0
		resp.Message = "success"
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", req.ResourceType)
	}
	return
}

func (l *FetchResourcesLogic) getAllTaskDefine(client scheduler.SchedulerClient, namespace string) ([]map[string]interface{}, error) {
	grpcResp, err := client.ListTaskDefines(l.ctx, &scheduler.ListTaskDefinesRequest{
		Metadata: map[string]string{
			"namespace": namespace,
		},
		PageSize: 100, // 设置适当的页面大小
		PageNum:  1,   // 从第一页开始
	})
	if err != nil {
		return nil, fmt.Errorf("call scheduler service failed: %v", err)
	}

	taskDefines := make([]map[string]interface{}, 0, len(grpcResp.Data))
	for _, taskDefine := range grpcResp.Data {
		taskDefines = append(taskDefines, map[string]interface{}{
			"spec":   taskDefine.Spec,
			"status": taskDefine.Status,
		})
	}
	return taskDefines, nil
}
