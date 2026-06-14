// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

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

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
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

	passwordHash, err := utils.HashPassword(req.Password, l.svcCtx.Config.Auth.PasswordPepper)
	if err != nil {
		return nil, utils.InternalError("hash password failed")
	}

	_, err = l.svcCtx.UserDAO.Insert(l.ctx, username, passwordHash)
	if err != nil {
		if utils.IsDuplicateEntryError(err, "uk_username") {
			return nil, utils.Conflict("username already exists")
		}
		return nil, utils.InternalError("create user failed: " + err.Error())
	}

	return &types.RegisterResponse{Registered: true}, nil
}
