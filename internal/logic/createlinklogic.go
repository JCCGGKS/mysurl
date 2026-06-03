// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLinkLogic {
	return &CreateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateLink validates the input URL, reuses an existing short link when the
// normalized URL already exists, otherwise generates a short code, inserts the
// mapping, and returns the final short-link payload.
func (l *CreateLinkLogic) CreateLink(req *types.CreateLinkRequest) (resp *types.CreateLinkResponse, err error) {
	if l.svcCtx.DB == nil {
		return nil, utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.ShortLinkDAO == nil {
		return nil, utils.InternalError("short link dao is not configured")
	}
	if l.svcCtx.CodeManager == nil {
		return nil, utils.InternalError("short code manager is not configured")
	}

	if err := utils.ValidateLongURL(req.LongURL); err != nil {
		return nil, utils.BadRequest(err.Error())
	}

	originalURL := strings.TrimSpace(req.LongURL)
	normalizedURL := utils.NormalizeOriginalURL(originalURL)
	urlHash := utils.HashOriginalURL(normalizedURL)

	bloomExists, err := l.svcCtx.ShortLinkCache.BloomExists(l.ctx, normalizedURL)
	if err != nil {
		l.Errorf("check normalized url bloom failed: %v", err)
	} else if bloomExists {
		l.Infof("create link bloom possible hit, normalized_url=%s", normalizedURL)
	} else {
		l.Infof("create link bloom miss, normalized_url=%s", normalizedURL)
	}

	if bloomExists {
		shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, normalizedURL)
		if cacheErr != nil {
			l.Errorf("get normalized url cache failed: %v", cacheErr)
		} else if shortCode != "" {
			l.Infof("create link hit long->short cache, normalized_url=%s short_code=%s", normalizedURL, shortCode)
			return l.buildCreateLinkResponse(shortCode, normalizedURL), nil
		}

		record, findErr := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, normalizedURL)
		if findErr != nil && !errors.Is(findErr, sqlx.ErrNotFound) {
			l.Errorf("query short link by original url failed: %v", findErr)
			return nil, utils.InternalError("query short link failed: " + findErr.Error())
		}
		if findErr == nil {
			l.fillCreateCaches(normalizedURL, record.ShortCode)
			l.Infof("create link hit mysql by original_url, normalized_url=%s short_code=%s", normalizedURL, record.ShortCode)
			return l.buildCreateLinkResponse(record.ShortCode, record.OriginalURL), nil
		}
	}

	result, fresh, err := l.svcCtx.FlightGroup.DoEx(createSingleflightKey(normalizedURL), func() (any, error) {
		shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, normalizedURL)
		if cacheErr != nil {
			l.Errorf("get normalized url cache in singleflight failed: %v", cacheErr)
		} else if shortCode != "" {
			return createLinkResult{
				shortCode:   shortCode,
				originalURL: normalizedURL,
				source:      "long->short cache",
			}, nil
		}

		record, findErr := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, normalizedURL)
		if findErr != nil && !errors.Is(findErr, sqlx.ErrNotFound) {
			return nil, findErr
		}
		if findErr == nil {
			l.fillCreateCaches(normalizedURL, record.ShortCode)
			return createLinkResult{
				shortCode:   record.ShortCode,
				originalURL: record.OriginalURL,
				source:      "mysql by original_url",
			}, nil
		}

		shortCode, genErr := l.svcCtx.CodeManager.GenerateShortCode(l.ctx, normalizedURL, urlHash)
		if genErr != nil {
			return nil, genErr
		}

		l.fillCreateCaches(normalizedURL, shortCode)
		return createLinkResult{
			shortCode:   shortCode,
			originalURL: normalizedURL,
			source:      fmt.Sprintf("new short code by provider=%s", l.svcCtx.CodeManager.Provider()),
		}, nil
	})
	if err != nil {
		l.Errorf("create link in singleflight failed: %v", err)
		return nil, utils.InternalError("create link failed: " + err.Error())
	}

	createResult, ok := result.(createLinkResult)
	if !ok {
		l.Errorf("invalid create link result type: %T", result)
		return nil, utils.InternalError("invalid create link result")
	}

	logSource := createResult.source
	if !fresh {
		logSource = fmt.Sprintf("%s via singleflight", logSource)
	}

	l.Infof("create link hit %s, normalized_url=%s short_code=%s", logSource, normalizedURL, createResult.shortCode)
	return l.buildCreateLinkResponse(createResult.shortCode, createResult.originalURL), nil
}

type createLinkResult struct {
	shortCode   string
	originalURL string
	source      string
}

func (l *CreateLinkLogic) fillCreateCaches(normalizedURL, shortCode string) {
	if err := l.svcCtx.ShortLinkCache.SetLongToShort(l.ctx, normalizedURL, shortCode); err != nil {
		l.Errorf("set normalized url cache failed: %v", err)
	}
	if err := l.svcCtx.ShortLinkCache.BloomAdd(l.ctx, normalizedURL); err != nil {
		l.Errorf("add normalized url bloom failed: %v", err)
	}
}

func (l *CreateLinkLogic) buildCreateLinkResponse(shortCode, originalURL string) *types.CreateLinkResponse {
	return &types.CreateLinkResponse{
		ShortCode:   shortCode,
		ShortURL:    utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, shortCode),
		OriginalURL: originalURL,
	}
}

func createSingleflightKey(normalizedURL string) string {
	return "create:" + normalizedURL
}
