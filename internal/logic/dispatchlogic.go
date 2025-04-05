package logic

import (
	"context"
	"encoding/json"
	"fmt"

	ie "entry/internal/error"
	"entry/internal/logic/utils"
	"entry/internal/proto/storage"
	"entry/internal/svc"
	"entry/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDispatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchLogic) Dispatch(req *types.DispatchResourceRequest) (resp *types.DispatchResourceResponse, err error) {
	resp = &types.DispatchResourceResponse{
		Code:    int(ie.Success),
		Message: ie.Success.GetMessage(),
		Extra:   make(map[string]interface{}),
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
		svcErr := ie.WrapError(ie.ServiceUnavailable, fmt.Errorf("failed to create storage RPC client"), "storage service initialization failed")
		resp.Code = int(ie.GetErrorCode(svcErr))
		resp.Message = svcErr.Error()
		return resp, svcErr
	}

	// metadata中没有指定namespace，指定default
	if _, ok := req.Metadata["namespace"]; !ok {
		req.Metadata["namespace"] = "default"
	}

	rpcResp, dispatchErr := l.dispatchResourceByStorage(storageClient, req)
	if dispatchErr != nil {
		l.Logger.Errorf("storage dispatch failed: %v", dispatchErr)
		resp.Code = int(ie.GetErrorCode(dispatchErr))
		resp.Message = dispatchErr.(*ie.Error).Message
		return resp, dispatchErr
	}

	resp.Code = int(ie.Success)
	resp.Message = fmt.Sprintf("%s Resource Successfully", req.Event)
	resp.Extra = rpcResp
	return
}

// func (l *DispatchLogic) dispatchResource(client scheduler.SchedulerClient, resourceType string, event string, spec map[string]interface{}, metadata map[string]interface{}) (map[string]interface{}, error) {
// 	switch resourceType {
// 	case "1": // ImageBuild
// 		return nil, nil
// 	case "2": // Sync
// 		return utils.DispatchSynchronizerEvent(l.ctx, client, event, spec, metadata)
// 	case "3": // ApiRuntime
// 		return nil, nil
// 	default:
// 		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
// 	}
// }

func (l *DispatchLogic) dispatchResourceByStorage(client storage.StorageClient, req *types.DispatchResourceRequest) (map[string]interface{}, error) {
	typeMapStr := `
	{
		"sync": 1,
		"apiruntime": 2,
		"imagebuild": 3
	}
	`
	var typeMap map[string]int
	err := json.Unmarshal([]byte(typeMapStr), &typeMap)
	if err != nil {
		logx.Errorf("Failed to unmarshal JSON: %v", err)
		return nil, ie.WrapError(ie.DataProcessError, err, "Failed to unmarshal resource type mapping")
	}

	storageUtil := utils.StorageSvcUtil{
		Client: client,
		Ctx:    l.ctx,
	}

	switch req.Event {
	case "createTestData":
		return storageUtil.CreateTaskData()
	case "updateTestData":
		return storageUtil.UpdateTaskData()
	case "deleteTestData":
		return storageUtil.DeleteTaskData()
	// case "getTestData":
	// 	return storageUtil.GetTaskData()
	// case "getAllTestData":
	// 	return storageUtil.GetTestDataList()
	case "createSceneConfig":
		return storageUtil.CreateSceniroConfig()
	case "updateSceneConfig":
		return storageUtil.UpdateSceniroConfig()
	case "deleteSceneConfig":
		return storageUtil.DeleteSceniroConfig()
	// case "getSceneConfig":
	// 	return storageUtil.GetSceniroConfig()
	// case "getAllSceneConfig":
	// 	return storageUtil.GetAllSceniroConfig()
	case "createTask":
		return storageUtil.CreateTask(req.Spec, req.Metadata, typeMap[req.ResourceType])
	case "updateTask":
		updateTaskId := req.Metadata["taskId"].(string)
		return storageUtil.UpdateTask(updateTaskId, req.Spec, req.Metadata, typeMap[req.ResourceType])
	case "deleteTask":
		updateTaskId := req.Metadata["taskId"].(string)
		return storageUtil.DeleteTask(updateTaskId)
	case "execTask":
		execTask := req.Metadata["taskId"].(string)
		return storageUtil.ExecuteTask(execTask)
	default:
		return nil, ie.NewError(ie.InvalidEventType, "unsupported event: %s", req.Event)
	}
}

// func (l *DispatchLogic) dispatchTaskDefineEvent(client scheduler.SchedulerClient, event string, spec, metadata map[string]interface{}) (map[string]interface{}, error) {
// 	// 构建 metadata
// 	newMetadata := make(map[string]string)
// 	for key, value := range metadata {
// 		newMetadata[key] = value.(string)
// 	}
// 	switch event {
// 	case "create":
// 		relatedImage := make(map[string]string)
// 		for key, value := range spec["related_image"].(map[string]interface{}) {
// 			relatedImage[key] = value.(string)
// 		}
// 		definitionBts, err := json.Marshal(spec["definition"])
// 		if err != nil {
// 			return nil, err
// 		}

// 		resp, err := client.CreateTaskDefine(l.ctx, &scheduler.CreateTaskDefineRequest{
// 			Metadata: newMetadata,
// 			Spec: &scheduler.TaskDefineSpec{
// 				IdlCode:    spec["idl_code"].(string),
// 				IdlType:    spec["idl_type"].(string),
// 				IdlName:    spec["idl_name"].(string),
// 				IdlVersion: spec["idl_version"].(string),
// 				RelatedImage: &scheduler.TaskDefineSpec_RelatedImage{
// 					Builder:   relatedImage["builder"],
// 					Digest:    relatedImage["digest"],
// 					Version:   relatedImage["version"],
// 					Namespace: relatedImage["namespace"],
// 				},
// 				Definition: string(definitionBts),
// 			},
// 		})

// 		logx.Error(resp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result, err := structToMap(resp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return result, nil

// 	case "update":
// 		newSpec := &scheduler.TaskDefineSpec{}
// 		if v, ok := spec["idl_code"]; ok {
// 			newSpec.IdlCode = v.(string)
// 		}
// 		if v, ok := spec["idl_name"]; ok {
// 			newSpec.IdlName = v.(string)
// 		}
// 		if v, ok := spec["idl_version"]; ok {
// 			newSpec.IdlVersion = v.(string)
// 		}
// 		if v, ok := spec["idl_type"]; ok {
// 			newSpec.IdlType = v.(string)
// 		}
// 		if v, ok := spec["related_image"]; ok {
// 			vmap := v.(map[string]interface{})
// 			newSpec.RelatedImage = &scheduler.TaskDefineSpec_RelatedImage{
// 				Builder:   vmap["builder"].(string),
// 				Digest:    vmap["digest"].(string),
// 				Version:   vmap["version"].(string),
// 				Namespace: vmap["namespace"].(string),
// 			}
// 		}
// 		if v, ok := spec["definition"]; ok {
// 			newDefinition, err := json.Marshal(v)
// 			if err != nil {
// 				return nil, err
// 			}
// 			newSpec.Definition = string(newDefinition)
// 		}

// 		logx.Info(newSpec)

// 		resp, err := client.UpdateTaskDefine(l.ctx, &scheduler.UpdateTaskDefineRequest{
// 			Spec:     newSpec,
// 			Metadata: newMetadata,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		result, err := structToMap(resp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return result, nil
// 	case "delete":
// 		resp, err := client.DeleteTaskDefine(l.ctx, &scheduler.DeleteTaskDefineRequest{
// 			Metadata: newMetadata,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		result, err := structToMap(resp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return result, nil

// 	default:
// 		{
// 			return nil, fmt.Errorf("unsupported event type: %s", event)
// 		}
// 	}

// 	return nil, errors.ErrUnsupported
// }

func structToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bts, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bts, &result); err != nil {
		return nil, err
	}
	return result, nil
}
