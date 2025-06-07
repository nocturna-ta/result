package controller

import (
	"context"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/infrastructures/custresp"
	"strconv"
	"time"
)

// GetVoteResult godoc
// @Summary Get vote result by ID
// @Description Get a specific vote result by its ID
// @Tags Results
// @Accept json
// @Produce json
// @Param id path string true "Vote Result ID"
// @Success 200 {object} jsonResponse{data=response.VoteResultResponse} "Vote result data"
// @Router /v1/results/votes/{id} [get]
func (api *API) GetVoteResult(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetVoteResult")
	defer span.End()

	id := req.Params("id")
	if id == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "vote result ID is required",
			Code:    400,
		})
	}

	result, err := api.voteResult.GetVoteResultByID(ctx, id)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(result), nil
}

// GetVoteResultByElectionPair godoc
// @Summary Get vote results by election pair ID
// @Description Get all vote results for a specific election pair
// @Tags Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Param limit query int false "Limit the number of results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} jsonResponse{data=[]response.VoteResultResponse} "List of vote results"
// @Router /v1/results/elections/{election_pair_id}/votes [get]
func (api *API) GetVoteResultByElectionPair(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetVoteResultByElectionPair")
	defer span.End()

	electionPairID := req.Params("election_pair_id")
	limit, err := strconv.Atoi(req.Query("limit", "50"))
	offset, err := strconv.Atoi(req.Query("offset", "0"))
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid limit or offset",
			Code:    400,
		})
	}

	results, err := api.voteResult.GetVoteResultsByElectionPair(ctx, electionPairID, limit, offset)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetVoteResultByRegion godoc
// @Summary Get vote results by region
// @Description Get all vote results for a specific region
// @Tags Results
// @Accept json
// @Produce json
// @Param region path string true "Region"
// @Param limit query int false "Limit the number of results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} jsonResponse{data=[]response.VoteResultResponse} "List of vote results"
// @Router /v1/results/regions/{region}/votes [get]
func (api *API) GetVoteResultByRegion(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetVoteResultByRegion")
	defer span.End()

	region := req.Params("region")
	limit, err := strconv.Atoi(req.Query("limit", "50"))
	offset, err := strconv.Atoi(req.Query("offset", "0"))
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid limit or offset",
			Code:    400,
		})
	}

	results, err := api.voteResult.GetVoteResultsByRegion(ctx, region, limit, offset)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetVoteResultByStatus godoc
// @Summary Get vote results by status
// @Description Get all vote results with a specific status
// @Tags Results
// @Accept json
// @Produce json
// @Param status query string true "Vote Result Status"
// @Param limit query int false "Limit the number of results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} jsonResponse{data=[]response.VoteResultResponse} "List of vote results"
// @Router /v1/results/votes [get]
func (api *API) GetVoteResultByStatus(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetVoteResultByStatus")
	defer span.End()

	status := req.Query("status")
	limit, err := strconv.Atoi(req.Query("limit", "50"))
	offset, err := strconv.Atoi(req.Query("offset", "0"))
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid limit or offset",
			Code:    400,
		})
	}

	results, err := api.voteResult.GetVoteResultsByStatus(ctx, status, limit, offset)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetElectionResults godoc
// @Summary Get election results by election pair ID
// @Description Get detailed election results for a specific election pair
// @Tags Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Success 200 {object} jsonResponse{data=response.ElectionVoteResultResponse} "Election results data"
// @Router /v1/results/elections/{election_pair_id} [get]
func (api *API) GetElectionResults(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetElectionResults")
	defer span.End()
	electionPairID := req.Params("election_pair_id")

	result, err := api.voteResult.GetElectionResults(ctx, electionPairID)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(result), nil
}

// GetElectionResultsByRegion godoc
// @Summary Get election results by region
// @Description Get detailed election results for a specific region
// @Tags Results
// @Accept json
// @Produce json
// @Param region path string true "Region"
// @Success 200 {object} jsonResponse{data=[]response.ElectionVoteResultResponse} "List of election results"
// @Router /v1/results/regions/{region}/elections [get]
func (api *API) GetElectionResultsByRegion(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetElectionResultsByRegion")
	defer span.End()

	region := req.Params("region")

	results, err := api.voteResult.GetElectionResultsByRegion(ctx, region)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetRegionResults godoc
// @Summary Get region results by region
// @Description Get detailed region results for a specific region
// @Tags Results
// @Accept json
// @Produce json
// @Param region path string true "Region"
// @Success 200 {object} jsonResponse{data=response.RegionVoteResultResponse} "Region results data"
// @Router /v1/results/regions/{region} [get]
func (api *API) GetRegionResults(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetRegionResults")
	defer span.End()

	region := req.Params("region")

	result, err := api.voteResult.GetRegionResults(ctx, region)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(result), nil
}

// GetRegionStatistics godoc
// @Summary Get region statistics
// @Description Get statistical data for all regions
// @Tags Results
// @Accept json
// @Produce json
// @Success 200 {object} jsonResponse{data=[]response.RegionVoteResultResponse} "List of region statistics"
// @Router /v1/results/regions [get]
func (api *API) GetRegionStatistics(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetRegionStatistics")
	defer span.End()

	results, err := api.voteResult.GetRegionStatistics(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetOverallStatistics godoc
// @Summary Get overall vote statistics
// @Description Get overall vote statistics including total votes, valid votes, and invalid votes
// @Tags Results
// @Accept json
// @Produce json
// @Success 200 {object} jsonResponse{data=response.VoteStatisticsResponse} "Overall vote statistics data"
// @Router /v1/results/statistics [get]
func (api *API) GetOverallStatistics(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetOverallStatistics")
	defer span.End()

	stats, err := api.voteResult.GetOverallStatistics(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(stats), nil
}

// GetDailyStatistics godoc
// @Summary Get daily vote statistics
// @Description Get daily vote statistics for a specified date range
// @Tags Results
// @Accept json
// @Produce json
// @Param start_date query string true "Start date in YYYY-MM-DD format"
// @Param end_date query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} jsonResponse{data=[]response.VoteStatisticsResponse} "List of daily vote statistics"
// @Router /v1/results/statistics/daily [get]
func (api *API) GetDailyStatistics(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.GetDailyStatistics")
	defer span.End()

	startDateStr := req.Query("start_date")
	endDateStr := req.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid start date format, expected YYYY-MM-DD",
			Code:    400,
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid end date format, expected YYYY-MM-DD",
			Code:    400,
		})
	}

	results, err := api.voteResult.GetDailyStatistics(ctx, startDate, endDate)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// CountVotesByStatus godoc
// @Summary Count votes by status
// @Description Count the number of votes with a specific status
// @Tags Results
// @Accept json
// @Produce json
// @Param status query string true "Vote Result Status"
// @Success 200 {object} jsonResponse{data=map[string]any} "Count of votes"
// @Router /v1/results/votes/count [get]
func (api *API) CountVotesByStatus(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.CountVotesByStatus")
	defer span.End()

	status := req.Query("status")

	count, err := api.voteResult.CountVotesByStatus(ctx, status)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]any{
		"count":  count,
		"status": status,
	}), nil

}

// CountVotesByElectionPair godoc
// @Summary Count votes by election pair ID
// @Description Count the number of votes for a specific election pair
// @Tags Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Success 200 {object} jsonResponse{data=map[string]any} "Count of votes"
// @Router /v1/results/elections/{election_pair_id}/count [get]
func (api *API) CountVotesByElectionPair(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.CountVotesByElectionPair")
	defer span.End()

	electionPairID := req.Params("election_pair_id")

	count, err := api.voteResult.CountVotesByElectionPair(ctx, electionPairID)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]any{
		"count": count,
	}), nil
}

// CountVotesByRegion godoc
// @Summary Count votes by region
// @Description Count the number of votes in a specific region
// @Tags Results
// @Accept json
// @Produce json
// @Param region path string true "Region"
// @Success 200 {object} jsonResponse{data=map[string]any} "Count of votes"
// @Router /v1/results/regions/{region}/count [get]
func (api *API) CountVotesByRegion(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultController.CountVotesByRegion")
	defer span.End()

	region := req.Params("region")

	count, err := api.voteResult.CountVotesByRegion(ctx, region)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]any{
		"count":  count,
		"region": region,
	}), nil
}
