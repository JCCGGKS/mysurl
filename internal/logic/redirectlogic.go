// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"strings"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Redirect validates the short code, loads the active short-link record,
// increments its visit count, and returns the target URL for the handler to
// issue the HTTP redirect.
func (l *RedirectLogic) Redirect(req *types.RedirectRequest) (string, error) {
	if l.svcCtx.DB == nil {
		return "", utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.ShortLinkDAO == nil {
		return "", utils.InternalError("short link dao is not configured")
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		return "", utils.BadRequest("code is required")
	}

	record, err := l.svcCtx.ShortLinkDAO.FindAvailableByCode(l.ctx, code)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return "", utils.NotFound("short link not found")
		}

		l.Errorf("query short link by code failed: %v", err)
		return "", utils.InternalError("query short link failed: "+ err.Error())
	}

	if err := l.svcCtx.ShortLinkDAO.IncrementVisitCount(l.ctx, record.ID); err != nil {
		l.Errorf("increment visit count failed: %v", err)
		return "", utils.InternalError("increment visit count failed")
	}

	return record.OriginalURL, nil
}
