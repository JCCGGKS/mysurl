package logic

import (
	"context"
	"strings"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserOperationLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserOperationLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserOperationLogsLogic {
	return &ListUserOperationLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserOperationLogsLogic) ListUserOperationLogs(req *types.ListUserOperationLogsRequest) (*types.UserOperationLogListResponse, error) {
	if l.svcCtx.UserOperationLogDAO == nil {
		return nil, utils.InternalError("user operation log dao is not configured")
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

	action := strings.TrimSpace(req.Action)
	result := strings.TrimSpace(req.Result)

	total, err := l.svcCtx.UserOperationLogDAO.CountByUserID(l.ctx, claims.UserID, action, result)
	if err != nil {
		l.Errorf("count user operation logs failed: %v", err)
		return nil, utils.InternalError("count user operation logs failed: " + err.Error())
	}

	records, err := l.svcCtx.UserOperationLogDAO.ListByUserIDWithCursor(l.ctx, claims.UserID, req.LastID, limit+1, action, result)
	if err != nil {
		l.Errorf("list user operation logs failed: %v", err)
		return nil, utils.InternalError("list user operation logs failed: " + err.Error())
	}

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}

	items := make([]types.UserOperationLogItem, 0, len(records))
	var nextLastID uint64
	for _, record := range records {
		item := types.UserOperationLogItem{
			ID:        record.ID,
			Action:    record.Action,
			Result:    record.Result,
			CreatedAt: record.CreatedAt.Unix(),
		}
		if record.Reason != nil {
			item.Reason = *record.Reason
		}

		items = append(items, item)
		nextLastID = record.ID
	}

	return &types.UserOperationLogListResponse{
		Items:      items,
		Total:      total,
		Limit:      limit,
		HasMore:    hasMore,
		NextLastID: nextLastID,
	}, nil
}
