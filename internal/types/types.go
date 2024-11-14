// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type DispatchRequest struct {
	TaskType     string `json:"taskType"`     // 任务类型
	ScheduleType string `json:"scheduleType"` // 指定要创建的crd类型
	Spec         string `json:"spec"`         // 任务的定义
	Verify       string `json:"verify"`       // 任务定义数据结构的验证
	Image        string `json:"image"`        // 关联的镜像
	Version      string `json:"version"`      // 版本号
}

type DispatchResponse struct {
	Code    int    `json:"code"`    // 响应状态码
	Message string `json:"message"` // 响应信息
}

type GetClusterCRsRequest struct {
	Type string `json:"type"` // 类型
}

type GetClusterCRsResponse struct {
	Code    int                      `json:"code"`    // 响应状态码
	Message string                   `json:"message"` // 响应信息
	Data    []map[string]interface{} `json:"data"`    // 数据
}

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}
