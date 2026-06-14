package logic

import (
	"context"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutRequest) error {
	if l.svcCtx.UserRefreshTokenDAO == nil {
		return utils.InternalError("user refresh token dao is not configured")
	}
	if req.RefreshToken == "" {
		return utils.Unauthorized("refresh token is required")
	}

	return l.svcCtx.UserRefreshTokenDAO.RevokeIfExists(l.ctx, utils.HashRefreshToken(req.RefreshToken))
}
