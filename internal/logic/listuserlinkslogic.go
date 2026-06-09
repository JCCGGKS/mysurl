// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLinksLogic {
	return &ListUserLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLinksLogic) ListUserLinks() (*types.UserLinkListResponse, error) {
	if l.svcCtx.ShortLinkDAO == nil {
		return nil, utils.InternalError("short link dao is not configured")
	}

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, utils.Unauthorized("authorization token is required")
	}

	records, err := l.svcCtx.ShortLinkDAO.ListByUserID(l.ctx, claims.UserID)
	if err != nil {
		l.Errorf("list user links failed: %v", err)
		return nil, utils.InternalError("list user links failed: " + err.Error())
	}

	items := make([]types.UserLinkItem, 0, len(records))
	for _, record := range records {
		items = append(items, types.UserLinkItem{
			ID:          record.ID,
			ShortCode:   record.ShortCode,
			ShortURL:    utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, record.ShortCode),
			OriginalURL: record.OriginalURL,
			VisitCount:  record.VisitCount,
			CreatedAt:   record.CreatedAt.Unix(),
		})
	}

	return &types.UserLinkListResponse{
		Items: items,
	}, nil
}
