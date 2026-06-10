package logic

import (
	"context"
	"strings"

	"mysurl1/internal/model"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

const batchCreateLinksLimit = 20

type BatchCreateLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchCreateLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateLinksLogic {
	return &BatchCreateLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateLinksLogic) BatchCreateLinks(req *types.BatchCreateLinksRequest) (*types.BatchCreateLinksResponse, error) {
	createLogic := NewCreateLinkLogic(l.ctx, l.svcCtx)
	if err := createLogic.ensureCreateReady(); err != nil {
		return nil, err
	}
	if len(req.LongURLs) == 0 {
		err := utils.BadRequest("long_urls is required")
		l.setBatchOperationLog(0, model.UserOperationResultFailed, err.Error())
		return nil, err
	}
	if len(req.LongURLs) > batchCreateLinksLimit {
		err := utils.BadRequest("long_urls exceeds limit 20")
		l.setBatchOperationLog(0, model.UserOperationResultFailed, err.Error())
		return nil, err
	}

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, utils.Unauthorized("authorization token is required")
	}

	items := make([]types.BatchCreateLinkItem, 0, len(req.LongURLs))
	successCache := make(map[string]*createLinkResult, len(req.LongURLs))
	successCount := 0

	for index, rawURL := range req.LongURLs {
		item := types.BatchCreateLinkItem{
			Index:   index,
			LongURL: rawURL,
		}

		if err := utils.ValidateLongURL(rawURL); err != nil {
			item.Error = err.Error()
			items = append(items, item)
			continue
		}

		normalizedURL := utils.NormalizeOriginalURL(strings.TrimSpace(rawURL))
		if cached := successCache[normalizedURL]; cached != nil {
			item.Success = true
			item.ShortCode = cached.ShortCode
			item.ShortURL = utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, cached.ShortCode)
			item.OriginalURL = cached.OriginalURL
			items = append(items, item)
			successCount++
			continue
		}

		result, err := createLogic.createOrReuseLink(claims.UserID, rawURL)
		if err != nil {
			item.Error = err.Error()
			items = append(items, item)
			continue
		}

		successCache[normalizedURL] = result
		item.Success = true
		item.ShortCode = result.ShortCode
		item.ShortURL = utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, result.ShortCode)
		item.OriginalURL = result.OriginalURL
		items = append(items, item)
		successCount++
	}

	if len(req.LongURLs)-successCount > 0 {
		l.setBatchOperationLog(claims.UserID, model.UserOperationResultFailed, firstBatchFailureReason(items))
	} else {
		l.setBatchOperationLog(claims.UserID, model.UserOperationResultSuccess, "")
	}

	return &types.BatchCreateLinksResponse{
		Items:        items,
		Total:        len(req.LongURLs),
		SuccessCount: successCount,
		FailedCount:  len(req.LongURLs) - successCount,
	}, nil
}

func (l *BatchCreateLinksLogic) setBatchOperationLog(userID uint64, result, reason string) {
	utils.SetOperationLogPayload(l.ctx, utils.OperationLogPayload{
		UserID: userID,
		Action: model.UserOperationActionCreateLinkBatch,
		Result: result,
		Reason: reason,
	})
}

func firstBatchFailureReason(items []types.BatchCreateLinkItem) string {
	for _, item := range items {
		if item.Error != "" {
			return item.Error
		}
	}

	return "batch create failed"
}
