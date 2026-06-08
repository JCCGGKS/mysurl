// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"mysurl1/internal/config"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.AuthResponse, err error) {
	authConf, err := ensureAuthConfig(l.svcCtx.Config.Auth)
	if err != nil {
		return nil, err
	}
	if l.svcCtx.DB == nil {
		return nil, utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.UserDAO == nil {
		return nil, utils.InternalError("user dao is not configured")
	}

	username := utils.NormalizeUsername(req.Username)
	if err := utils.ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if req.Password != req.ConfirmPassword {
		return nil, utils.BadRequest("password and confirm_password do not match")
	}

	_, err = l.svcCtx.UserDAO.FindByUsername(l.ctx, username)
	if err == nil {
		return nil, utils.Conflict("username already exists")
	}
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, utils.InternalError("query user failed: " + err.Error())
	}

	passwordHash, err := utils.HashPassword(req.Password, authConf.PasswordPepper)
	if err != nil {
		return nil, utils.InternalError("hash password failed")
	}

	userID, err := l.svcCtx.UserDAO.Insert(l.ctx, username, passwordHash)
	if err != nil {
		if utils.IsDuplicateEntryError(err, "uk_username") {
			return nil, utils.Conflict("username already exists")
		}
		return nil, utils.InternalError("create user failed: " + err.Error())
	}

	return utils.BuildAuthResponse(authConf, utils.AuthClaims{
		UserID:   userID,
		Username: username,
	})
}

func ensureAuthConfig(auth config.AuthConf) (config.AuthConf, error) {
	if auth.JWTSecret == "" {
		return auth, utils.InternalError("auth jwt secret is not configured")
	}
	if auth.ExpireSeconds <= 0 {
		auth.ExpireSeconds = 86400
	}

	return auth, nil
}
