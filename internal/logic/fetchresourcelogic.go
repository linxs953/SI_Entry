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

// func (l *FetchResourceLogic) getTaskDefine(client scheduler.SchedulerClient, namespace, name string) (map[string]interface{}, error) {
// 	if namespace == "" {
// 		namespace = "default"
// 	}
// 	if name == "" {
// 		return nil, errors.New("name is empty")
// 	}
// 	grpcResp, err := client.GetTaskDefine(l.ctx, &scheduler.GetTaskDefineRequest{
// 		Metadata: map[string]string{
// 			"name":      name,
// 			"namespace": namespace,
// 		},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("call scheduler service failed: %v", err)
// 	}

// 	if grpcResp.Data == nil {
// 		return nil, fmt.Errorf("TaskDefine obj is nil")
// 	}

// 	data := map[string]interface{}{
// 		"spec":   grpcResp.Data.Spec,
// 		"status": grpcResp.Data.Status,
// 	}
// 	return data, nil
// }

func (l *FetchResourceLogic) FetchResource(req *types.FetchResourceRequest) (resp *types.FetchResourceResponse, err error) {
	resp = &types.FetchResourceResponse{}
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
		svcErr := ie.WrapError(ie.ServiceUnavailable, fmt.Errorf("failed to create storage RPC client"), "storage service initialization failed")
		resp.Code = int(ie.GetErrorCode(svcErr))
		resp.Message = svcErr.Error()
		return resp, svcErr
	}
	return l.fetchOne(storageClient, req)
}

func (l *FetchResourceLogic) fetchOne(client storage.StorageClient, req *types.FetchResourceRequest) (resp *types.FetchResourceResponse, err error) {
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
	case "getTask":
		{
			taskId := req.Metadata["taskId"].(string)
			rpcResp, err := storageUtil.GetTask(taskId)
			if err != nil {
				return nil, err
			}

			logx.Error(rpcResp)

			// 处理 rpcResp，提取 Header 和 Spec
			header, ok := rpcResp["header"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid response format: header not found")
			}

			// spec, ok := rpcResp["Spec"].(map[string]interface{})
			// if !ok {
			// 	spec = make(map[string]interface{})
			// 	// return nil, fmt.Errorf("invalid response format: spec not found")
			// }

			if _, ok := header["code"].(float64); !ok {
				return nil, fmt.Errorf("rpc response header.code invalid")
			}

			if _, ok := header["message"].(string); !ok {
				return nil, fmt.Errorf("rpc response header.msg invalid")
			}

			data := make(map[string]interface{})
			data["spec"] = rpcResp["Spec"]
			data["meta"] = rpcResp["meta"]
			data["type"] = rpcResp["type"]
			data["create_at"] = rpcResp["create_at"]
			data["update_at"] = rpcResp["update_at"]
			resp = &types.FetchResourceResponse{
				Code:    int(header["code"].(float64)), // 假设 code 是 float64 类型
				Message: header["message"].(string),    // 假设 message 是 string 类型
				Data:    data,
			}
			return resp, nil
		}
	}
	return nil, fmt.Errorf("event %s not supported", req.Event)
}
