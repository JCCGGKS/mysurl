package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mysurl1/internal/logic"
	types "mysurl1/internal/schema"
	"mysurl1/internal/svc"
	"mysurl1/internal/utils"
)

func ListUserOperationLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListUserOperationLogsRequest
		if err := httpx.Parse(r, &req); err != nil {
			utils.WriteJSONError(w, r, utils.BadRequest(err.Error()))
			return
		}

		l := logic.NewListUserOperationLogsLogic(r.Context(), svcCtx)
		resp, err := l.ListUserOperationLogs(&req)
		if err != nil {
			utils.WriteJSONError(w, r, err)
		} else {
			utils.WriteJSONSuccess(w, r, resp)
		}
	}
}
