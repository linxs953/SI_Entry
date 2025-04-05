package utils

import (
	"context"
	"encoding/json"
	errors "entry/internal/error"
	ie "entry/internal/error"
	"entry/internal/proto/storage"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageSvcUtil struct {
	Client storage.StorageClient
	Ctx    context.Context
}

func (l *StorageSvcUtil) CreateTask(spec, metadata map[string]interface{}, resourceType int) (map[string]interface{}, error) {
	client := l.Client

	// Extract name and description from metadata
	var name, desc string
	if v, ok := metadata["name"]; ok {
		name = v.(string)
	}
	if v, ok := metadata["desc"]; ok {
		desc = v.(string)
	}

	// Create request with basic info
	req := &storage.CreateTaskRequest{
		Name: name,
		Type: storage.TaskType(resourceType),
		Desc: desc,
	}

	// Convert spec to appropriate task spec type
	specBytes, err := json.Marshal(spec)
	if err != nil {
		return nil, ie.WrapError(ie.DataSerializationError, err, "failed to serialize task spec")
	}

	// Create appropriate spec based on task type
	switch storage.TaskType(resourceType) {
	case 2:
		apiSpec := &storage.TaskAPISpec{}
		if err := json.Unmarshal(specBytes, apiSpec); err != nil {
			return nil, ie.WrapError(ie.DataSerializationError, err, "failed to parse API task spec")
		}
		req.Spec = &storage.CreateTaskRequest_ApiSpec{ApiSpec: apiSpec}

	case 1:
		syncSpec := &storage.TaskSyncSpec{}
		if err := json.Unmarshal(specBytes, syncSpec); err != nil {
			return nil, ie.WrapError(ie.DataSerializationError, err, "failed to parse Sync task spec")
		}
		req.Spec = &storage.CreateTaskRequest_SyncSpec{SyncSpec: syncSpec}

	default:
		return nil, ie.NewError(ie.InvalidTaskType, "unsupported task type: %d", resourceType)
	}

	// Handle task creation
	resp, err := client.CreateTask(l.Ctx, req)
	if err != nil {
		return nil, ie.WrapError(ie.TaskCreationFailed, err, "failed to create task")
	}
	if resp == nil {
		return nil, ie.NewError(ie.TaskCreationFailed, "received nil response from storage service")
	}

	if resp.Header.Code == int64(errors.RPCSuccess) {
		resp.Header.Code = int64(errors.Success)
	}

	if resp.Header.Code != int64(errors.Success) {
		return nil, ie.NewError(ie.TaskCreationFailed, "task creation failed: %s", resp.Header.Message)
	}

	// Convert response to map
	respMap := make(map[string]interface{})
	taskBts, err := json.Marshal(resp.Meta)
	if err == nil {
		err = json.Unmarshal(taskBts, &respMap)
	}
	return respMap, err
}

func (l *StorageSvcUtil) GetAllTask() (map[string]interface{}, error) {
	client := l.Client

	resp, err := client.ListTasks(l.Ctx, &storage.Empty{})
	if err != nil {
		return nil, ie.WrapError(ie.TaskFetchListFailed, err, "获取任务列表失败")
	}
	// if resp == nil {
	// 	return nil, ie.NewError(ie.TaskFetchListFailed, "received nil response from storage service")
	// }
	if resp == nil {
		return nil, ie.NewError(ie.TaskFetchListFailed, "获取任务列表失败")
	}
	if resp.Header == nil {
		return nil, ie.NewError(ie.TaskFetchListFailed, "获取任务列表失败")
	}

	if resp.Header.Code == int64(errors.RPCSuccess) {
		resp.Header.Code = int64(errors.Success)
	}
	if resp.Header.Code != int64(errors.Success) {
		return nil, ie.NewError(ie.TaskFetchListFailed, "获取任务列表失败")
	}

	// Convert response to map
	respMap := make(map[string]interface{})

	// taskBts, err := json.Marshal(resp)
	// if err == nil {
	// 	err = json.Unmarshal(taskBts, &respMap)
	// }

	return respMap, err
}

func (l *StorageSvcUtil) GetTask(taskId string) (map[string]interface{}, error) {
	client := l.Client
	req := storage.GetTaskRequest{
		TaskId: taskId,
	}

	resp, err := client.GetTask(l.Ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task with id=[%s]: %v", taskId, err)
	}

	respMap := make(map[string]interface{})

	taskBts, err := json.Marshal(resp)
	if err == nil {
		err = json.Unmarshal(taskBts, &respMap)
	}
	return respMap, err
}

func (l *StorageSvcUtil) UpdateTask(taskId string, spec, meta map[string]interface{}, resourceType int) (map[string]interface{}, error) {
	client := l.Client

	// Extract name and description from metadata
	var name, desc string
	if v, ok := meta["name"]; ok {
		name = v.(string)
	}
	if v, ok := meta["desc"]; ok {
		desc = v.(string)
	}

	// Create request with basic info
	req := &storage.UpdateTaskRequest{
		TaskId: taskId,
		Name:   name,
		Desc:   desc,
	}

	// Convert spec to appropriate task spec type
	specBytes, err := json.Marshal(spec)
	if err != nil {
		return nil, ie.WrapError(ie.DataSerializationError, err, "failed to serialize task spec")
	}

	// Create appropriate spec based on task type
	switch storage.TaskType(resourceType) {
	case 2:
		apiSpec := &storage.TaskAPISpec{}
		if err := json.Unmarshal(specBytes, apiSpec); err != nil {
			return nil, ie.WrapError(ie.DataSerializationError, err, "failed to parse API task spec")
		}
		req.Spec = &storage.UpdateTaskRequest_ApiSpec{ApiSpec: apiSpec}

	case 1:
		syncSpec := &storage.TaskSyncSpec{}
		if err := json.Unmarshal(specBytes, syncSpec); err != nil {
			return nil, ie.WrapError(ie.DataSerializationError, err, "failed to parse Sync task spec")
		}
		req.Spec = &storage.UpdateTaskRequest_SyncSpec{SyncSpec: syncSpec}

	default:
		return nil, ie.NewError(ie.InvalidTaskType, "unsupported task type: %d", resourceType)
	}

	// Handle task update
	resp, err := client.UpdateTask(l.Ctx, req)
	if err != nil {
		return nil, ie.WrapError(ie.TaskUpdateFailed, err, "failed to update task")
	}
	if resp == nil {
		return nil, ie.NewError(ie.TaskUpdateFailed, "received nil response from storage service")
	}
	if resp.Header.Code == int64(errors.RPCSuccess) {
		resp.Header.Code = int64(errors.Success)
	}
	if resp.Header.Code != int64(errors.Success) {
		logx.Error(resp.Header.Code, errors.Success)
		return nil, ie.NewError(ie.TaskUpdateFailed, "task update failed: %s", resp.Header.Message)
	}

	// Convert response to map
	respMap := make(map[string]interface{})
	respMap["taskId"] = resp.Meta.TaskId
	respMap["name"] = resp.Meta.TaskName
	respMap["desc"] = resp.Meta.TaskDesc
	respMap["spec"] = resp.Spec
	respMap["updateTime"] = resp.UpdateAt
	// taskBts, err := json.Marshal(resp)
	// if err == nil {
	// 	err = json.Unmarshal(taskBts, &respMap)
	// }
	return respMap, err
}

func (l *StorageSvcUtil) DeleteTask(taskId string) (map[string]interface{}, error) {
	client := l.Client
	req := &storage.DeleteTaskRequest{
		TaskId: taskId,
	}

	resp, err := client.DeleteTask(l.Ctx, req)
	if err != nil {
		return nil, ie.WrapError(ie.TaskDeletionFailed, err, "failed to delete task")
	}
	if resp == nil {
		return nil, ie.NewError(ie.TaskDeletionFailed, "received nil response from storage service")
	}

	if resp.Header.Code == int64(errors.RPCSuccess) {
		resp.Header.Code = int64(errors.Success)
	}

	if resp.Header.Code != int64(errors.Success) {
		return nil, ie.NewError(ie.TaskDeletionFailed, "task deletion failed: %s", resp.Header.Message)
	}

	// 返回成功删除的任务ID
	return map[string]interface{}{"taskId": taskId}, nil
}

func (l *StorageSvcUtil) ExecuteTask(taskId string) (map[string]interface{}, error) {
	client := l.Client
	resp, err := client.ExecuteTask(l.Ctx, &storage.ExecuteTaskRequest{TaskId: taskId})
	logx.Error(err)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ie.NewError(ie.TaskExecutionFailed, "received nil response from storage service")
	}
	if resp.Header.Code == int64(errors.RPCSuccess) {
		resp.Header.Code = int64(errors.Success)
	}
	if resp.Header.Code != int64(errors.Success) {
		return nil, ie.NewError(ie.TaskExecutionFailed, "task execution failed: %s", resp.Header.Message)
	}
	return map[string]interface{}{"taskId": taskId}, nil
}

// 测试数据
func (l *StorageSvcUtil) CreateTaskData() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) UpdateTaskData() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) DeleteTaskData() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) GetTaskData() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) GetTestDataList() (map[string]interface{}, error) {
	return nil, nil
}

// 自动化场景配置
func (l *StorageSvcUtil) CreateSceniroConfig() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) UpdateSceniroConfig() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) DeleteSceniroConfig() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) GetSceniroConfig() (map[string]interface{}, error) {
	return nil, nil
}

func (l *StorageSvcUtil) GetAllSceniroConfig() (map[string]interface{}, error) {
	return nil, nil
}
