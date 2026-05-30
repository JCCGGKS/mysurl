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

	candidates, err := l.svcCtx.ShortLinkDAO.FindAvailableByHash(l.ctx, urlHash)
	if err != nil {
		l.Errorf("query short links by hash failed: %v", err)
		return nil, utils.InternalError("query short links failed")
	}

	for _, candidate := range candidates {
		if utils.NormalizeOriginalURL(candidate.OriginalURL) == normalizedURL {
			return l.buildCreateLinkResponse(candidate.ShortCode, candidate.OriginalURL), nil
		}
	}

	shortCode, genErr := l.svcCtx.CodeManager.GenerateShortCode(l.ctx, originalURL, urlHash)
	if genErr != nil {
		l.Errorf("generate short code failed: %v", genErr)
		return nil, utils.InternalError("generate short code failed")
	}

	return l.buildCreateLinkResponse(shortCode, originalURL), nil
}

func (l *CreateLinkLogic) buildCreateLinkResponse(shortCode, originalURL string) *types.CreateLinkResponse {
	return &types.CreateLinkResponse{
		ShortCode:   shortCode,
		ShortURL:    utils.BuildShortURL(l.svcCtx.Config.Short.BaseURL, shortCode),
		OriginalURL: originalURL,
	}
}
