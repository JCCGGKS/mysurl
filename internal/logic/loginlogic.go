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
		l.setLoginOperationLog(0, model.UserOperationResultFailed, err.Error())
		return nil, err
	}
	if req.Password == "" {
		err := utils.BadRequest("password is required")
		l.setLoginOperationLog(0, model.UserOperationResultFailed, err.Error())
		return nil, err
	}

	user, err := l.svcCtx.UserDAO.FindByUsername(l.ctx, username)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			err := utils.Unauthorized("username or password is invalid")
			l.setLoginOperationLog(0, model.UserOperationResultFailed, err.Error())
			return nil, err
		}
		err := utils.InternalError("query user failed: " + err.Error())
		l.setLoginOperationLog(0, model.UserOperationResultFailed, err.Error())
		return nil, err
	}

	if err := utils.ComparePassword(user.PasswordHash, req.Password, authConf.PasswordPepper); err != nil {
		authErr := utils.Unauthorized("username or password is invalid")
		l.setLoginOperationLog(user.ID, model.UserOperationResultFailed, authErr.Error())
		return nil, authErr
	}

	l.setLoginOperationLog(user.ID, model.UserOperationResultSuccess, "")

	return utils.BuildAuthResponse(authConf, utils.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (l *LoginLogic) setLoginOperationLog(userID uint64, result, reason string) {
	utils.SetOperationLogPayload(l.ctx, utils.OperationLogPayload{
		UserID: userID,
		Action: model.UserOperationActionLogin,
		Result: result,
		Reason: reason,
	})
}
