// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strings"

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

func (l *ListUserLinksLogic) ListUserLinks(req *types.ListUserLinksRequest) (*types.UserLinkListResponse, error) {
	if l.svcCtx.ShortLinkDAO == nil {
		return nil, utils.InternalError("short link dao is not configured")
	}

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, utils.Unauthorized("authorization token is required")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = utils.DefaultPageSize
	}
	if limit > utils.MaxPageSize {
		limit = utils.MaxPageSize
	}

	shortCode := strings.TrimSpace(req.ShortCode)
	originalURL := strings.TrimSpace(req.OriginalURL)
	total, err := l.svcCtx.ShortLinkDAO.CountByUserID(l.ctx, claims.UserID, shortCode, originalURL)
	if err != nil {
		l.Errorf("count user links failed: %v", err)
		return nil, utils.InternalError("count user links failed: " + err.Error())
	}

	records, err := l.svcCtx.ShortLinkDAO.ListByUserIDWithCursor(
		l.ctx,
		claims.UserID,
		shortCode,
		originalURL,
		req.LastID,
		limit+1, // 探测是否还有下一页
	)
	if err != nil {
		l.Errorf("list user links failed: %v", err)
		return nil, utils.InternalError("list user links failed: " + err.Error())
	}

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}

	items := make([]types.UserLinkItem, 0, len(records))
	var nextLastID uint64
	for _, record := range records {
		items = append(items, types.UserLinkItem{
			ID:          record.ID,
			ShortCode:   record.ShortCode,
			ShortURL:    utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, record.ShortCode),
			OriginalURL: record.OriginalURL,
			VisitCount:  record.VisitCount,
			CreatedAt:   record.CreatedAt.Unix(),
		})
		nextLastID = record.ID
	}

	return &types.UserLinkListResponse{
		Items:      items,
		Total:      total,
		Limit:      limit,
		HasMore:    hasMore,
		NextLastID: nextLastID,
	}, nil
}
