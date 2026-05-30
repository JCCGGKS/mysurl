// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
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
			utils.WriteRedirectError(w, utils.BadRequest(err.Error()))
			return
		}

		l := logic.NewRedirectLogic(r.Context(), svcCtx)
		targetURL, err := l.Redirect(&req)
		if err != nil {
			utils.WriteRedirectError(w, err)
		} else {
			http.Redirect(w, r, targetURL, http.StatusFound)
		}
	}
}
