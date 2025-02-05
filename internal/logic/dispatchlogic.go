package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"entry/internal/logic/utils"
	"entry/internal/proto/scheduler"
	"entry/internal/svc"
	"entry/internal/types"
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
		Code:    0,
		Message: "success",
	}

	// 获取gRPC连接
	conn, err := l.svcCtx.ClientManager.GetClient("scheduler")
	if err != nil {
		l.Logger.Errorf("scheduler service not found, continuing without scheduler")
		resp.Code = -1
		resp.Message = err.Error()
		return resp, err
	}

	// 创建scheduler client
	client := scheduler.NewSchedulerClient(conn)
	if client == nil {
		l.Logger.Errorf("scheduler client creation failed, continuing without scheduler")
		resp.Code = 1
		resp.Message = "scheduler client creation failed, continuing without scheduler"
		return resp, nil
	}

	// metadata中没有指定namespace，指定default
	if _, ok := req.Metadata["namespace"]; !ok {
		req.Metadata["namespace"] = "default"
	}

	// 尝试调用scheduler服务
	_, err = l.dispatchResource(client, req.ResourceType, req.Event, req.Spec, req.Metadata)
	if err != nil {
		// 记录错误但不中断流程
		l.Logger.Errorf("scheduler dispatch failed: %v, continuing without scheduler", err)
		resp.Code = 2
		resp.Message = "scheduler dispatch failed"
		return resp, nil
	}
	resp.Code = 0
	resp.Message = fmt.Sprintf("%s Resource Successfully", req.Event)
	return
}

func (l *DispatchLogic) dispatchResource(client scheduler.SchedulerClient, resourceType string, event string, spec map[string]interface{}, metadata map[string]interface{}) (map[string]interface{}, error) {
	switch resourceType {
	case "1": // ImageBuild
		return nil, nil
	case "2": // Sync
		return utils.DispatchSynchronizerEvent(l.ctx, client, event, spec, metadata)
	case "3": // ApiRuntime
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
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
