package controller

import (
	"context"
	"encoding/json"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/infrastructures/custresp"
	request2 "github.com/nocturna-ta/result/internal/usecases/request"
)

// InsertVoteResult godoc
// @Summary Insert Vote Result
// @Description Insert a new vote result
// @Tags VoteResult
// @Accept json
// @Produce json
// @Param request body request.VoteResultEntry true "Vote Result Entry"
// @Success 200 {object} jsonResponse{data=response.EntryResponse} "Vote Result Entry"
// @Route POST /vote-result
func (api *API) InsertVoteResult(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultController.InsertVoteResult")
	defer span.End()

	var request request2.VoteResultEntry
	err := json.Unmarshal(req.RawBody(), &request)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.voteResult.InsertVoteResult(ctx, request)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}
