package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	ie "entry/internal/error"
	"entry/internal/logic/utils"
	"entry/internal/proto/storage"
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
	resp = &types.FetchAllResourcesResponse{
		Code:    int(ie.Success),
		Message: "success",
		Data:    nil,
	}

	conn, err := l.svcCtx.ClientManager.GetClient("storage")
	if err != nil {
		l.Logger.Errorf("storage service not found: %v", err)
		svcErr := ie.WrapError(ie.ServiceUnavailable, err, "storage service not found")
		resp.Code = int(ie.GetErrorCode(svcErr))
		resp.Message = svcErr.Error()
		return resp, svcErr
	}

	storageClient := storage.NewStorageClient(conn)
	if storageClient == nil {
		l.Logger.Error("storage rpc client creation failed")
		err := fmt.Errorf("failed to create storage RPC client")
		svcErr := ie.WrapError(ie.ServiceUnavailable, err, "storage service initialization failed")
		resp.Code = int(ie.GetErrorCode(svcErr))
		resp.Message = svcErr.Error()
		return resp, svcErr
	}

	return l.fetchAll(storageClient, req)
}

// func (l *FetchResourcesLogic) getAllTaskDefine(client scheduler.SchedulerClient, namespace string) ([]map[string]interface{}, error) {
// 	grpcResp, err := client.ListTaskDefines(l.ctx, &scheduler.ListTaskDefinesRequest{
// 		Metadata: map[string]string{
// 			"namespace": namespace,
// 		},
// 		PageSize: 100, // 设置适当的页面大小
// 		PageNum:  1,   // 从第一页开始
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("call scheduler service failed: %v", err)
// 	}

// 	taskDefines := make([]map[string]interface{}, 0, len(grpcResp.Data))
// 	for _, taskDefine := range grpcResp.Data {
// 		taskDefines = append(taskDefines, map[string]interface{}{
// 			"spec":   taskDefine.Spec,
// 			"status": taskDefine.Status,
// 		})
// 	}
// 	return taskDefines, nil
// }

func (l *FetchResourcesLogic) fetchAll(client storage.StorageClient, req *types.FetchAllResourcesRequest) (resp *types.FetchAllResourcesResponse, err error) {
	typeMapStr := `
		{
			"sync": 1,
			"apiruntime": 2,
			"imagebuild": 3
		}
		`
	var typeMap map[string]int
	err = json.Unmarshal([]byte(typeMapStr), &typeMap)
	if err != nil {
		logx.Errorf("Failed to unmarshal JSON: %v", err)
		return nil, ie.WrapError(ie.DataProcessError, err, "Failed to unmarshal resource type mapping")
	}

	storageUtil := utils.StorageSvcUtil{
		Client: client,
		Ctx:    l.ctx,
	}
	switch req.Event {
	case "getAllTask":
		{
			taskList, err := storageUtil.GetAllTask()
			if err != nil {
				return nil, err
			}

			logx.Error(taskList["data"])
			data, err := json.Marshal(taskList["data"])
			if err != nil {
				return nil, err
			}
			var dataMap []map[string]interface{}
			err = json.Unmarshal(data, &dataMap)
			if err != nil {
				return nil, err
			}

			return &types.FetchAllResourcesResponse{
				Code:    int(ie.Success),
				Message: "获取任务列表成功",
				Data:    dataMap,
			}, nil
		}
	}
	return nil, fmt.Errorf("event %s not supported", req.Event)
}
