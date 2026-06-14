// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"time"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.AuthResponse, err error) {
	authConf, err := utils.EnsureAuthConfig(l.svcCtx.Config.Auth)
	if err != nil {
		return nil, err
	}
	if l.svcCtx.DB == nil {
		return nil, utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.UserDAO == nil {
		return nil, utils.InternalError("user dao is not configured")
	}
	if l.svcCtx.UserRefreshTokenDAO == nil {
		return nil, utils.InternalError("user refresh token dao is not configured")
	}

	username := utils.NormalizeUsername(req.Username)
	if err := utils.ValidateUsername(username); err != nil {
		return nil, err
	}
	if req.Password == "" {
		return nil, utils.BadRequest("password is required")
	}

	user, err := l.svcCtx.UserDAO.FindByUsername(l.ctx, username)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, utils.Unauthorized("username or password is invalid")
		}
		return nil, utils.InternalError("query user failed: " + err.Error())
	}

	if err := utils.ComparePassword(user.PasswordHash, req.Password, authConf.PasswordPepper); err != nil {
		return nil, utils.Unauthorized("username or password is invalid")
	}

	authResp, err := utils.BuildAuthResponse(authConf, utils.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
	})
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.UserRefreshTokenDAO.Insert(
		l.ctx,
		user.ID,
		utils.HashRefreshToken(authResp.RefreshToken),
		time.Unix(authResp.RefreshExpiresAt, 0),
	); err != nil {
		return nil, utils.InternalError("save refresh token failed: " + err.Error())
	}

	return authResp, nil
}
