package handler

import (
	"net/http"

	"mysurl1/internal/logic"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogoutRequest
		if err := httpx.Parse(r, &req); err != nil {
			utils.WriteJSONError(w, r, utils.BadRequest(err.Error()))
			return
		}

		l := logic.NewLogoutLogic(r.Context(), svcCtx)
		if err := l.Logout(&req); err != nil {
			utils.WriteJSONError(w, r, err)
			return
		}

		utils.WriteJSONSuccess(w, r, map[string]any{"success": true})
	}
}
