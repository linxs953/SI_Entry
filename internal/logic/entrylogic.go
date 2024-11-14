package logic

import (
	"context"

	"entry/internal/svc"
	"entry/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntryLogic {
	return &EntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntryLogic) Entry(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
