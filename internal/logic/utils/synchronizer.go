package utils

import (
	"context"
	"entry/internal/proto/scheduler"
	"fmt"
)

// 对接scheduler，调用Synchronizer的方法，操作crd对象
// create / update / delete

func DispatchSynchronizerEvent(ctx context.Context, client scheduler.SchedulerClient, event string, spec, metadata map[string]interface{}) (map[string]interface{}, error) {
	newMetadata := make(map[string]string)
	for key, value := range metadata {
		newMetadata[key] = value.(string)
	}
	switch event {
	case "create":
		return createSynchronizer(ctx, client, spec)

	case "update":
		return updateSynchronizer(ctx, client, spec)

	case "delete":
		return deleteSynchronizer(ctx, client, metadata)

	default:
		{
			return nil, fmt.Errorf("unsupported event type: %s", event)
		}
	}

}

// todo： 这几个方法待更新scheduler的proto文件后，再补充
func createSynchronizer(ctx context.Context, client scheduler.SchedulerClient, spec map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil

}

func updateSynchronizer(ctx context.Context, client scheduler.SchedulerClient, spec map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func deleteSynchronizer(ctx context.Context, client scheduler.SchedulerClient, metadata map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
