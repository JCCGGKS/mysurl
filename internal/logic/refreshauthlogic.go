package logic

import (
	"context"
	"errors"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RefreshAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAuthLogic {
	return &RefreshAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshAuthLogic) Refresh(req *types.RefreshRequest) (*types.RefreshResponse, error) {
	authConf, err := utils.EnsureAuthConfig(l.svcCtx.Config.Auth)
	if err != nil {
		return nil, err
	}
	if l.svcCtx.UserRefreshTokenDAO == nil {
		return nil, utils.InternalError("user refresh token dao is not configured")
	}

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return nil, utils.Unauthorized("refresh token is required")
	}

	tokenHash := utils.HashRefreshToken(refreshToken)
	record, err := l.svcCtx.UserRefreshTokenDAO.FindActiveByTokenHash(l.ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, utils.Unauthorized("refresh token is invalid")
		}
		return nil, utils.InternalError("query refresh token failed: " + err.Error())
	}

	user, err := l.svcCtx.UserDAO.FindByID(l.ctx, record.UserID)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, utils.Unauthorized("refresh token is invalid")
		}
		return nil, utils.InternalError("query user failed: " + err.Error())
	}

	tokenPair, err := utils.CreateTokenPair(authConf, utils.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
	})
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.UserRefreshTokenDAO.Rotate(
		l.ctx,
		tokenHash,
		record.UserID,
		tokenPair.RefreshTokenHash,
		tokenPair.RefreshExpiresAt,
	); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, utils.Unauthorized("refresh token is invalid")
		}
		return nil, utils.InternalError("rotate refresh token failed: " + err.Error())
	}

	return &types.RefreshResponse{
		AccessToken:      tokenPair.AccessToken,
		AccessExpiresAt:  tokenPair.AccessExpiresAt.Unix(),
		RefreshToken:     tokenPair.RefreshToken,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt.Unix(),
	}, nil
}
