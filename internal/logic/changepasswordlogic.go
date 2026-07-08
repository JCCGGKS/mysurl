// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"mysurl1/internal/model"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) error {
	if l.svcCtx.DB == nil {
		return utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.UserDAO == nil {
		return utils.InternalError("user dao is not configured")
	}

	var user *model.User
	var err error

	// Try to get user ID from JWT token first
	claims, ok := utils.GetAuthClaims(l.ctx)
	if ok && claims != nil {
		user, err = l.svcCtx.UserDAO.FindByID(l.ctx, claims.UserID)
	} else if req.Username != "" {
		// If no token, use username from request
		user, err = l.svcCtx.UserDAO.FindByUsername(l.ctx, req.Username)
	} else {
		return utils.BadRequest("username is required when not logged in")
	}

	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return utils.NotFound("user not found")
		}
		return utils.InternalError("query user failed: " + err.Error())
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}
	if req.NewPassword != req.ConfirmPassword {
		return utils.BadRequest("new_password and confirm_password do not match")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword, l.svcCtx.Config.Auth.PasswordPepper)
	if err != nil {
		return utils.InternalError("hash password failed")
	}

	if err := l.svcCtx.UserDAO.UpdatePassword(l.ctx, user.ID, passwordHash); err != nil {
		return utils.InternalError("update password failed: " + err.Error())
	}

	return nil
}
