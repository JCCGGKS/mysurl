// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mysurl1/internal/logic"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"
)

func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRedirectLogic(r.Context(), svcCtx)
		targetURL, err := l.Redirect(&req)
		if err != nil {
			var httpErr *utils.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Message, httpErr.StatusCode)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, targetURL, http.StatusFound)
		}
	}
}
