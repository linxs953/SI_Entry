package errors

// ErrorCode represents a unique error code in the system
type ErrorCode int

// Error codes for different categories
const (
	// Success
	Success ErrorCode = 0

	// Client Errors (1000-1999)
	InvalidParameter     ErrorCode = 1000
	ValidationFailed     ErrorCode = 1001
	ResourceNotFound     ErrorCode = 1002
	ResourceDuplicate    ErrorCode = 1003
	Unauthorized         ErrorCode = 1004
	Forbidden            ErrorCode = 1005
	RequestTimeout       ErrorCode = 1006
	TooManyRequests      ErrorCode = 1007
	InvalidResourceType  ErrorCode = 1008
	InvalidEventType     ErrorCode = 1009
	InvalidMetadata      ErrorCode = 1010
	InvalidConfiguration ErrorCode = 1011

	// Server Errors (2000-2999)
	InternalError          ErrorCode = 2000
	DatabaseError          ErrorCode = 2001
	CacheError             ErrorCode = 2002
	NetworkError           ErrorCode = 2003
	ServiceUnavailable     ErrorCode = 2004
	ConfigError            ErrorCode = 2005
	StorageServiceError    ErrorCode = 2006
	SchedulerServiceError  ErrorCode = 2007
	RPCError               ErrorCode = 2008
	ClientCreationError    ErrorCode = 2009
	ServiceConnectionError ErrorCode = 2010
	DataSerializationError ErrorCode = 2011
	ServiceTimeoutError    ErrorCode = 2012

	// Business Logic Errors (3000-3999)
	TaskExecutionFailed   ErrorCode = 3000
	InvalidTaskStatus     ErrorCode = 3001
	SceneConfigInvalid    ErrorCode = 3002
	ResourceBusy          ErrorCode = 3003
	DataProcessError      ErrorCode = 3004
	TaskCreationFailed    ErrorCode = 3005
	TaskUpdateFailed      ErrorCode = 3006
	TaskDeletionFailed    ErrorCode = 3007
	TaskNotFound          ErrorCode = 3008
	InvalidTaskType       ErrorCode = 3009
	TaskValidationFailed  ErrorCode = 3010
	SceneCreationFailed   ErrorCode = 3011
	SceneUpdateFailed     ErrorCode = 3012
	SceneDeletionFailed   ErrorCode = 3013
	SceneNotFound         ErrorCode = 3014
	InvalidSceneType      ErrorCode = 3015
	SceneValidationFailed ErrorCode = 3016
	TaskFetchListFailed   ErrorCode = 3017

	// rpc error
	RPCSuccess ErrorCode = 100
)

// GetMessage returns a default message for error codes
func (code ErrorCode) GetMessage() string {
	switch code {
	case Success:
		return "success"
	case InvalidParameter:
		return "invalid parameter"
	case ValidationFailed:
		return "validation failed"
	case ResourceNotFound:
		return "resource not found"
	case ResourceDuplicate:
		return "resource already exists"
	case Unauthorized:
		return "unauthorized"
	case Forbidden:
		return "forbidden"
	case RequestTimeout:
		return "request timeout"
	case TooManyRequests:
		return "too many requests"
	case InternalError:
		return "internal server error"
	case DatabaseError:
		return "database error"
	case CacheError:
		return "cache error"
	case NetworkError:
		return "network error"
	case ServiceUnavailable:
		return "service unavailable"
	case StorageServiceError:
		return "storage service error"
	case SchedulerServiceError:
		return "scheduler service error"
	case RPCError:
		return "RPC communication error"
	case ClientCreationError:
		return "client creation failed"
	case ServiceConnectionError:
		return "service connection failed"
	case DataSerializationError:
		return "data serialization error"
	case ServiceTimeoutError:
		return "service timeout error"
	case InvalidResourceType:
		return "invalid resource type"
	case InvalidEventType:
		return "invalid event type"
	case InvalidMetadata:
		return "invalid metadata"
	case InvalidConfiguration:
		return "invalid configuration"
	case TaskCreationFailed:
		return "task creation failed"
	case TaskUpdateFailed:
		return "task update failed"
	case TaskDeletionFailed:
		return "task deletion failed"
	case TaskNotFound:
		return "task not found"
	case TaskFetchListFailed:
		return "task list fetch error"
	case InvalidTaskType:
		return "invalid task type"
	case TaskValidationFailed:
		return "task validation failed"
	case SceneCreationFailed:
		return "scene creation failed"
	case SceneUpdateFailed:
		return "scene update failed"
	case SceneDeletionFailed:
		return "scene deletion failed"
	case SceneNotFound:
		return "scene not found"
	case InvalidSceneType:
		return "invalid scene type"
	case SceneValidationFailed:
		return "scene validation failed"
	case ConfigError:
		return "configuration error"
	case TaskExecutionFailed:
		return "task execution failed"
	case InvalidTaskStatus:
		return "invalid task status"
	case SceneConfigInvalid:
		return "scene configuration invalid"
	case ResourceBusy:
		return "resource is busy"
	case DataProcessError:
		return "data processing error"

	default:
		return "unknown error"
	}
}
