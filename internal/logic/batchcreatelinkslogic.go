package logic

import (
	"context"
	"strings"

	codestrategy "mysurl1/internal/logic/code_strategy"
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
	return l.batchCreateLinks(req, false)
}

func (l *BatchCreateLinksLogic) batchCreateLinks(req *types.BatchCreateLinksRequest, retried bool) (*types.BatchCreateLinksResponse, error) {
	createLogic := NewCreateLinkLogic(l.ctx, l.svcCtx)
	if err := createLogic.ensureCreateReady(); err != nil {
		return nil, err
	}
	if len(req.LongURLs) == 0 {
		return nil, utils.BadRequest("long_urls is required")
	}
	if len(req.LongURLs) > batchCreateLinksLimit {
		return nil, utils.BadRequest("long_urls exceeds limit 20")
	}

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, utils.Unauthorized("authorization token is required")
	}

	items := make([]types.BatchCreateLinkItem, 0, len(req.LongURLs))
	uniqueURLs := make([]string, 0, len(req.LongURLs))
	seenURLs := make(map[string]struct{}, len(req.LongURLs))
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
		if _, exists := seenURLs[normalizedURL]; !exists {
			uniqueURLs = append(uniqueURLs, normalizedURL)
			seenURLs[normalizedURL] = struct{}{}
		}
		item.OriginalURL = normalizedURL
		items = append(items, item)
	}

	existingResults, err := l.loadExistingBatchResults(claims.UserID, uniqueURLs)
	if err != nil {
		return nil, utils.InternalError("find existing links failed: " + err.Error())
	}

	pendingRecords := make([]model.ShortLink, 0, len(uniqueURLs))
	pendingURLs := make([]string, 0, len(uniqueURLs))
	resultByURL := make(map[string]*createLinkResult, len(uniqueURLs))

	for _, normalizedURL := range uniqueURLs {
		if existing, ok := existingResults[normalizedURL]; ok {
			resultByURL[normalizedURL] = existing
			continue
		}

		pendingURLs = append(pendingURLs, normalizedURL)
		pendingRecords = append(pendingRecords, model.ShortLink{
			OriginalURL: normalizedURL,
			URLHash:     utils.HashOriginalURL(normalizedURL),
		})
	}

	if len(pendingRecords) > 0 {
		createdResults, err := l.batchCreateNewLinks(claims.UserID, pendingRecords)
		if err != nil {
			if utils.IsDuplicateEntryError(err, "uk_user_original_url") && !retried {
				return l.batchCreateLinks(req, true)
			}
			if utils.IsDuplicateEntryError(err, "uk_short_code") {
				return nil, utils.InternalError("batch create links failed: short code conflict")
			}
			return nil, utils.InternalError("batch create links failed: " + err.Error())
		}
		for i, created := range createdResults {
			resultByURL[pendingURLs[i]] = &createLinkResult{
				ShortCode:   created.ShortCode,
				OriginalURL: created.OriginalURL,
				Source:      createLinkSourceCreated,
			}
		}
	}

	for i := range items {
		if items[i].OriginalURL == "" {
			continue
		}

		result, ok := resultByURL[items[i].OriginalURL]
		if !ok {
			items[i].Error = "create link failed"
			continue
		}

		items[i].Success = true
		items[i].ShortCode = result.ShortCode
		items[i].ShortURL = utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, result.ShortCode)
		items[i].OriginalURL = result.OriginalURL
		successCount++
	}

	return &types.BatchCreateLinksResponse{
		Items:        items,
		Total:        len(req.LongURLs),
		SuccessCount: successCount,
		FailedCount:  len(req.LongURLs) - successCount,
	}, nil
}

func (l *BatchCreateLinksLogic) loadExistingBatchResults(userID uint64, normalizedURLs []string) (map[string]*createLinkResult, error) {
	results := make(map[string]*createLinkResult, len(normalizedURLs))
	if len(normalizedURLs) == 0 {
		return results, nil
	}

	missedURLs := make([]string, 0, len(normalizedURLs))
	for _, normalizedURL := range normalizedURLs {
		shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, userID, normalizedURL)
		if cacheErr == nil && shortCode != "" {
			results[normalizedURL] = &createLinkResult{
				ShortCode:   shortCode,
				OriginalURL: normalizedURL,
				Source:      createLinkSourceCacheHit,
			}
			continue
		}
		missedURLs = append(missedURLs, normalizedURL)
	}
	if len(missedURLs) == 0 {
		return results, nil
	}

	records, err := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURLs(l.ctx, userID, missedURLs)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		l.newCreateLinkLogic().fillCreateCaches(userID, record.OriginalURL, record.ShortCode)
		results[record.OriginalURL] = &createLinkResult{
			ShortCode:   record.ShortCode,
			OriginalURL: record.OriginalURL,
			Source:      createLinkSourceDBHit,
		}
	}

	return results, nil
}

func (l *BatchCreateLinksLogic) batchCreateNewLinks(userID uint64, records []model.ShortLink) ([]model.ShortLink, error) {
	if len(records) == 0 {
		return nil, nil
	}

	createLogic := l.newCreateLinkLogic()
	if createLogic.shortCodeProvider() == codestrategy.ProviderMySQLAutoIncrement {
		created, err := l.svcCtx.ShortLinkDAO.BatchCreateWithAutoIncrement(l.ctx, &userID, records)
		if err != nil {
			return nil, err
		}
		for _, record := range created {
			createLogic.fillCreateCaches(userID, record.OriginalURL, record.ShortCode)
		}
		return created, nil
	}

	for i := range records {
		shortCode, err := l.svcCtx.CodeManager.NextCode(l.ctx, createLogic.shortCodeProvider())
		if err != nil {
			return nil, err
		}
		records[i].ShortCode = shortCode
	}

	if err := l.svcCtx.ShortLinkDAO.BatchInsert(l.ctx, &userID, records); err != nil {
		return nil, err
	}
	for _, record := range records {
		createLogic.fillCreateCaches(userID, record.OriginalURL, record.ShortCode)
	}

	return records, nil
}

func (l *BatchCreateLinksLogic) newCreateLinkLogic() *CreateLinkLogic {
	return NewCreateLinkLogic(l.ctx, l.svcCtx)
}
