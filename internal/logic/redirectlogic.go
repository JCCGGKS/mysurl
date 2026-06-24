// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"
	"strings"

	"mysurl1/internal/dao"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
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

	cacheValue, cacheErr := l.svcCtx.ShortLinkCache.GetShortToLong(l.ctx, code)
	if cacheErr == nil && cacheValue != nil {
		return l.returnRedirectTarget(code, cacheValue.ID, cacheValue.OriginalURL, "short->long cache")
	}

	l.Errorf("get short code cache miss or failed: %v", cacheErr)

	bloomExists, bloomErr := l.svcCtx.ShortLinkCache.ShortCodeBloomExists(l.ctx, code)
	if bloomErr == nil && !bloomExists {
		return "", utils.NotFound("short link not found")
	}
	l.Errorf("check short code bloom hit or failed: %v", bloomErr)

	result, fresh, err := l.svcCtx.FlightGroup.DoEx(redirectSingleflightKey(code), func() (any, error) {
		record, err := l.svcCtx.ShortLinkDAO.FindAvailableByCode(l.ctx, code)
		if err != nil || record == nil {
			return nil, fmt.Errorf("look database miss or failed: %v", err)
		}

		if err := l.svcCtx.ShortLinkCache.ShortCodeBloomAdd(l.ctx, code); err != nil {
			l.Errorf("add short code bloom failed: %v", err)
		}
		if err := l.svcCtx.ShortLinkCache.SetShortToLong(l.ctx, code, dao.ShortToLongCacheValue{
			ID:          record.ID,
			OriginalURL: record.OriginalURL,
		}); err != nil {
			l.Errorf("set short code cache failed: %v", err)
		}

		return redirectLookupResult{
			id:          record.ID,
			originalURL: record.OriginalURL,
			source:      "mysql by short_code",
		}, nil
	})
	if err != nil {
		l.Errorf("query short link by code not found or failed: %v", err)
		return "", utils.InternalError("query short link failed: " + err.Error())
	}

	lookupResult, ok := result.(redirectLookupResult)
	if !ok {
		l.Errorf("invalid redirect lookup result type: %T", result)
		return "", utils.InternalError("invalid redirect lookup result")
	}

	logSource := lookupResult.source
	// fresh 表示当前请求是否是真正执行函数
	if !fresh {
		logSource = fmt.Sprintf("%s via singleflight", logSource)
	}

	return l.returnRedirectTarget(code, lookupResult.id, lookupResult.originalURL, logSource)
}

type redirectLookupResult struct {
	id          uint64
	originalURL string
	source      string
}

func (l *RedirectLogic) returnRedirectTarget(code string, id uint64, targetURL, source string) (string, error) {
	baseCount := uint64(0)
	if _, exists, err := l.svcCtx.ShortLinkCache.GetVisitCount(l.ctx, id); err != nil {
		l.Errorf("get visit count cache failed: %v", err)
	} else if !exists {
		dbCount, dbErr := l.svcCtx.VisitStatDAO.GetVisitCount(l.ctx, id)
		if dbErr != nil {
			l.Errorf("load visit base count failed: %v", dbErr)
		} else {
			baseCount = dbCount
		}
	}

	if err := l.svcCtx.ShortLinkCache.IncrVisitCount(l.ctx, id, baseCount); err != nil {
		l.Errorf("incr visit count failed: %v", err)
	}

	l.Infof("redirect hit %s, short_code=%s target=%s", source, code, targetURL)
	return targetURL, nil
}

func redirectSingleflightKey(code string) string {
	return "redirect:" + code
}
