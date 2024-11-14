package logic

import (
	"context"

	"entry/internal/svc"
	"entry/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClusterCRsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterCRsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterCRsLogic {
	return &GetClusterCRsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClusterCRsLogic) GetClusterCRs(req *types.GetClusterCRsRequest) (resp *types.GetClusterCRsResponse, err error) {
	resp = &types.GetClusterCRsResponse{
		Code:    0,
		Message: "success",
		Data: []map[string]interface{}{
			{
				"type":          req.Type,
				"schedule_type": "syncapi",
				"idlName":       "test-idl",
				"idlDesc":       "test-idl-desc",
				"version":       "v1",
				"image":         "test-image:latest",
				"spec":          map[string]interface{}{},
			},
			{
				"type":          req.Type,
				"schedule_type": "syncapi",
				"idlName":       "test-idl2",
				"idlDesc":       "test-idl-desc2",
				"version":       "v2",
				"image":         "test-image2:latest",
				"spec":          map[string]interface{}{},
			},
		},
	}
	return
}
