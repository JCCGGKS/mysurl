// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"strings"

	"mysurl1/internal/model"
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

type createLinkResult struct {
	RecordID    uint64
	ShortCode   string
	OriginalURL string
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

	claims, ok := utils.GetAuthClaims(l.ctx)
	if !ok || claims.UserID == 0 {
		return nil, utils.Unauthorized("authorization token is required")
	}
	result, err := l.createOrReuseLink(claims.UserID, req.LongURL)
	if err != nil {
		return nil, err
	}

	l.setCreateOperationLog(claims.UserID, result.RecordID, result.ShortCode)
	return l.buildCreateLinkResponse(result.ShortCode, result.OriginalURL), nil
}

func (l *CreateLinkLogic) createOrReuseLink(userID uint64, longURL string) (*createLinkResult, error) {
	originalURL := strings.TrimSpace(longURL)
	normalizedURL := utils.NormalizeOriginalURL(originalURL)
	urlHash := utils.HashOriginalURL(normalizedURL)

	shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, userID, normalizedURL)
	if cacheErr != nil {
		l.Errorf("get normalized url cache failed: %v", cacheErr)
	} else if shortCode != "" {
		l.Infof("create link hit long->short cache, user_id=%d normalized_url=%s short_code=%s", userID, normalizedURL, shortCode)
		record, err := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
		if err == nil {
			return &createLinkResult{
				RecordID:    record.ID,
				ShortCode:   record.ShortCode,
				OriginalURL: record.OriginalURL,
			}, nil
		}
	}

	record, err := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		l.Errorf("query short link by original url failed: %v", err)
		return nil, utils.InternalError("query short link failed: " + err.Error())
	}
	if err == nil {
		l.fillCreateCaches(userID, normalizedURL, record.ShortCode)
		l.Infof("create link hit mysql by user_id+original_url, user_id=%d normalized_url=%s short_code=%s", userID, normalizedURL, record.ShortCode)
		return &createLinkResult{
			RecordID:    record.ID,
			ShortCode:   record.ShortCode,
			OriginalURL: record.OriginalURL,
		}, nil
	}

	shortCode, err = l.createNewLink(userID, normalizedURL, urlHash)
	if err != nil {
		if utils.IsDuplicateEntryError(err, "uk_user_original_url") {
			record, findErr := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
			if findErr != nil {
				l.Errorf("query short link after duplicate original url failed: %v", findErr)
				return nil, utils.InternalError("query short link failed: " + findErr.Error())
			}

			l.fillCreateCaches(userID, normalizedURL, record.ShortCode)
			l.Infof("create link reuse after duplicate user_id+original_url, user_id=%d normalized_url=%s short_code=%s", userID, normalizedURL, record.ShortCode)
			return &createLinkResult{
				RecordID:    record.ID,
				ShortCode:   record.ShortCode,
				OriginalURL: record.OriginalURL,
			}, nil
		}

		l.Errorf("create new short link failed: %v", err)
		return nil, utils.InternalError("create link failed: " + err.Error())
	}

	l.fillCreateCaches(userID, normalizedURL, shortCode)
	l.Infof("create link generated new short code, provider=%s user_id=%d normalized_url=%s short_code=%s", l.svcCtx.CodeManager.Provider(), userID, normalizedURL, shortCode)
	record, err = l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
	if err != nil {
		l.Errorf("query short link after create failed: %v", err)
		return nil, utils.InternalError("query short link failed: " + err.Error())
	}

	return &createLinkResult{
		RecordID:    record.ID,
		ShortCode:   record.ShortCode,
		OriginalURL: record.OriginalURL,
	}, nil
}

func (l *CreateLinkLogic) createNewLink(userID uint64, normalizedURL, urlHash string) (string, error) {
	if l.svcCtx.CodeManager.IsMySQLAutoIncrement() {
		return l.svcCtx.ShortLinkDAO.CreateWithAutoIncrement(l.ctx, &userID, normalizedURL, urlHash)
	}

	shortCode, err := l.svcCtx.CodeManager.NextCode(l.ctx)
	if err != nil {
		return "", err
	}
	if err := l.svcCtx.ShortLinkDAO.Insert(l.ctx, &userID, shortCode, normalizedURL, urlHash); err != nil {
		return "", err
	}

	return shortCode, nil
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

func (l *CreateLinkLogic) setCreateOperationLog(userID, targetID uint64, shortCode string) {
	targetIDCopy := targetID
	targetCodeCopy := shortCode
	utils.SetOperationLogPayload(l.ctx, utils.OperationLogPayload{
		UserID:     userID,
		Action:     model.UserOperationActionCreateLink,
		Result:     model.UserOperationResultSuccess,
		TargetID:   &targetIDCopy,
		TargetCode: &targetCodeCopy,
	})
}
