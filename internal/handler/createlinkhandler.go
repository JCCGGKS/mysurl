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

func CreateLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateLinkLogic(r.Context(), svcCtx)
		resp, err := l.CreateLink(&req)
		if err != nil {
			var httpErr *utils.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Message, httpErr.StatusCode)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
