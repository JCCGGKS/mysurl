package handler

import (
	"net/http"

	"mysurl1/internal/logic"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshRequest
		if err := httpx.Parse(r, &req); err != nil {
			utils.WriteJSONError(w, r, utils.BadRequest(err.Error()))
			return
		}

		l := logic.NewRefreshAuthLogic(r.Context(), svcCtx)
		resp, err := l.Refresh(&req)
		if err != nil {
			utils.WriteJSONError(w, r, err)
			return
		}

		utils.WriteJSONSuccess(w, r, resp)
	}
}
