// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"mysurl1/internal/svc"
	"mysurl1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLinkLogic {
	return &CreateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLinkLogic) CreateLink(req *types.CreateLinkRequest) (resp *types.CreateLinkResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
