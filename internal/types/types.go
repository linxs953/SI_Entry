// Code generated by goctl. DO NOT EDIT.
package types

type FetchResourceRequest struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"name"`
	Namespace    string `json:"namespace,optional"`
}

type FetchResourceResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type FetchAllResourcesRequest struct {
	ResourceType string `json:"resource_type"`
	Namespace    string `json:"namespace,optional"`
}

type FetchAllResourcesResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type DispatchResourceRequest struct {
	ResourceType string                 `json:"type"`
	Event        string                 `json:"event"`
	Spec         map[string]interface{} `json:"spec"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type DispatchResourceResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
