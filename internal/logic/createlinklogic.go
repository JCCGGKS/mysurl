// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strings"

	codestrategy "mysurl1/internal/logic/code_strategy"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type createLinkResult struct {
	ShortCode   string
	OriginalURL string
	Source      string
}

const (
	createLinkSourceCacheHit = "cache_hit"
	createLinkSourceDBHit    = "database_hit"
	createLinkSourceCreated  = "created"
)

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
func (l *CreateLinkLogic) CreateLink(req *types.CreateLinkRequest) (resp *types.CreateLinkResponse, extData any, err error) {
	if err := l.ensureCreateReady(); err != nil {
		return nil, nil, err
	}

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, nil, utils.Unauthorized("authorization token is required")
	}
	if err := utils.ValidateLongURL(req.LongURL); err != nil {
		return nil, nil, utils.BadRequest(err.Error())
	}
	result, err := l.createOrReuseLink(claims.UserID, req.LongURL)
	if err != nil {
		return nil, nil, err
	}

	return l.buildCreateLinkResponse(result.ShortCode, result.OriginalURL), result.Source, nil
}

func (l *CreateLinkLogic) ensureCreateReady() error {
	if l.svcCtx.DB == nil {
		return utils.InternalError("mysql is not configured")
	}
	if l.svcCtx.ShortLinkDAO == nil {
		return utils.InternalError("short link dao is not configured")
	}
	if l.svcCtx.CodeManager == nil {
		return utils.InternalError("short code manager is not configured")
	}

	return nil
}

func (l *CreateLinkLogic) createOrReuseLink(userID uint64, longURL string) (*createLinkResult, error) {
	normalizedURL := utils.NormalizeOriginalURL(strings.TrimSpace(longURL))
	urlHash := utils.HashOriginalURL(normalizedURL)

	shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, userID, normalizedURL)
	if cacheErr == nil && shortCode != "" {
		l.Infof("create link hit long->short cache, user_id=%d normalized_url=%s short_code=%s", userID, normalizedURL, shortCode)
		return &createLinkResult{
			ShortCode:   shortCode,
			OriginalURL: normalizedURL,
			Source:      createLinkSourceCacheHit,
		}, nil
	}

	record, err := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
	if err == nil && record != nil {
		l.fillCreateCaches(userID, normalizedURL, record.ShortCode)
		l.Infof("create link hit mysql by user_id+original_url, user_id=%d normalized_url=%s short_code=%s", userID, normalizedURL, record.ShortCode)
		return &createLinkResult{
			ShortCode:   record.ShortCode,
			OriginalURL: record.OriginalURL,
			Source:      createLinkSourceDBHit,
		}, nil
	}

	shortCode, err = l.createNewLink(userID, normalizedURL, urlHash)
	if err != nil {
		l.Errorf("create new short link failed: %v", err)
		return nil, utils.InternalError("create link failed: " + err.Error())
	}

	l.fillCreateCaches(userID, normalizedURL, shortCode)
	l.Infof("create link generated new short code, provider=%s user_id=%d normalized_url=%s short_code=%s", l.shortCodeProvider(), userID, normalizedURL, shortCode)
	return &createLinkResult{
		ShortCode:   shortCode,
		OriginalURL: normalizedURL,
		Source:      createLinkSourceCreated,
	}, nil
}

func (l *CreateLinkLogic) createNewLink(userID uint64, normalizedURL, urlHash string) (string, error) {
	if l.shortCodeProvider() == codestrategy.ProviderMySQLAutoIncrement {
		return l.svcCtx.ShortLinkDAO.CreateWithAutoIncrement(l.ctx, &userID, normalizedURL, urlHash)
	}

	shortCode, err := l.svcCtx.CodeManager.NextCode(l.ctx, l.shortCodeProvider())
	if err != nil {
		return "", err
	}
	if err := l.svcCtx.ShortLinkDAO.Insert(l.ctx, &userID, shortCode, normalizedURL, urlHash); err != nil {
		return "", err
	}

	return shortCode, nil
}

func (l *CreateLinkLogic) shortCodeProvider() string {
	provider := strings.TrimSpace(l.svcCtx.Config.Short.Provider)
	if provider == "" {
		return codestrategy.ProviderMySQLAutoIncrement
	}

	return provider
}

func (l *CreateLinkLogic) fillCreateCaches(userID uint64, normalizedURL, shortCode string) {
	if err := l.svcCtx.ShortLinkCache.SetLongToShort(l.ctx, userID, normalizedURL, shortCode); err != nil {
		l.Errorf("set normalized url cache failed: %v", err)
	}
	if err := l.svcCtx.ShortLinkCache.ShortCodeBloomAdd(l.ctx, shortCode); err != nil {
		l.Errorf("add short code bloom failed: %v", err)
	}
}

func (l *CreateLinkLogic) buildCreateLinkResponse(shortCode, originalURL string) *types.CreateLinkResponse {
	return &types.CreateLinkResponse{
		ShortCode:   shortCode,
		ShortURL:    utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, shortCode),
		OriginalURL: originalURL,
	}
}
