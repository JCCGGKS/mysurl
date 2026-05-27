// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"mysurl1/internal/svc"
	"mysurl1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Mysurl1Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMysurl1Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Mysurl1Logic {
	return &Mysurl1Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Mysurl1Logic) Mysurl1(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
