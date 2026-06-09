// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"mysurl1/internal/logic"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"
)

func ListUserLinksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListUserLinksLogic(r.Context(), svcCtx)
		resp, err := l.ListUserLinks()
		if err != nil {
			utils.WriteJSONError(w, r, err)
		} else {
			utils.WriteJSONSuccess(w, r, resp)
		}
	}
}
