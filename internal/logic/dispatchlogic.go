package logic

import (
	"context"

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

func (l *DispatchLogic) Dispatch(req *types.DispatchRequest) (resp *types.DispatchResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
