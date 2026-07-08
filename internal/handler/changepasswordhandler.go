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

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			utils.WriteJSONError(w, r, utils.BadRequest(err.Error()))
			return
		}

		l := logic.NewChangePasswordLogic(r.Context(), svcCtx)
		err := l.ChangePassword(&req)
		if err != nil {
			utils.WriteJSONError(w, r, err)
		} else {
			utils.WriteJSONSuccess(w, r, nil)
		}
	}
}
