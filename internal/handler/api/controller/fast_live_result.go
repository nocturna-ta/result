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

// GetFastLiveElectionResults godoc
// @Summary Get fast live election results with percentages
// @Description Get real-time election results with percentages optimized for high performance using materialized views and Redis cache
// @Tags Fast Live Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Param include_cities query bool false "Include top cities data" default(true)
// @Param cities_limit query int false "Limit number of cities returned" default(10)
// @Success 200 {object} jsonResponse{data=response.FastElectionResultsResponse} "Fast live election results"
// @Router /v1/live/elections/{election_pair_id} [get]
func (api *API) GetFastLiveElectionResults(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetFastLiveElectionResults")
	defer span.End()

	electionPairID := req.Params("election_pair_id")
	if electionPairID == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
		})
	}

	// Get fast live results with cache
	results, err := api.fastLiveResultUc.GetLiveElectionResultsWithCache(ctx, electionPairID)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetFastCityResults godoc
// @Summary Get fast live city results
// @Description Get real-time city voting results optimized for high performance using materialized views and Redis cache
// @Tags Fast Live Results
// @Accept json
// @Produce json
// @Param city_name path string true "City Name"
// @Success 200 {object} jsonResponse{data=response.FastCityResultsResponse} "Fast live city results"
// @Router /v1/live/cities/{city_name} [get]
func (api *API) GetFastCityResults(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetFastCityResults")
	defer span.End()

	cityName := req.Params("city_name")
	if cityName == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "city name is required",
			Code:    400,
		})
	}

	// Get fast city results with cache
	results, err := api.fastLiveResultUc.GetLiveCityResultsWithCache(ctx, cityName)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	// Check if city has any results
	if len(results.ElectionResults) == 0 {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "no results found for city",
			Code:    404,
		})
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetFastElectionSummaries godoc
// @Summary Get fast summaries for all elections
// @Description Get real-time summaries for all active elections optimized for dashboards using materialized views and Redis cache
// @Tags Fast Live Results
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of elections returned" default(20)
// @Success 200 {object} jsonResponse{data=response.FastAllElectionsSummaryResponse} "Fast election summaries"
// @Router /v1/live/elections [get]
func (api *API) GetFastElectionSummaries(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetFastElectionSummaries")
	defer span.End()

	// Get all election summaries with cache
	results, err := api.fastLiveResultUc.GetAllElectionsSummaryWithCache(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	// Apply limit if requested
	limit, _ := strconv.Atoi(req.Query("limit", "20"))
	if limit > 0 && limit < len(results.Elections) {
		results.Elections = results.Elections[:limit]
		results.TotalElections = uint64(limit)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetFastCityRankings godoc
// @Summary Get fast city rankings for an election
// @Description Get real-time city rankings for a specific election optimized for high performance using materialized views and Redis cache
// @Tags Fast Live Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Param limit query int false "Limit number of cities returned" default(20)
// @Success 200 {object} jsonResponse{data=response.FastCityRankingsResponse} "Fast city rankings"
// @Router /v1/live/rankings/{election_pair_id} [get]
func (api *API) GetFastCityRankings(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetFastCityRankings")
	defer span.End()

	electionPairID := req.Params("election_pair_id")
	if electionPairID == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
		})
	}

	limit, _ := strconv.Atoi(req.Query("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Get fast city rankings with cache
	results, err := api.fastLiveResultUc.GetCityRankingsWithCache(ctx, electionPairID, limit)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// GetFastElectionSummary godoc
// @Summary Get fast election summary
// @Description Get real-time summary for a specific election optimized for high performance using materialized views and Redis cache
// @Tags Fast Live Results
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Success 200 {object} jsonResponse{data=response.FastElectionSummaryResponse} "Fast election summary"
// @Router /v1/live/elections/{election_pair_id}/summary [get]
func (api *API) GetFastElectionSummary(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetFastElectionSummary")
	defer span.End()

	electionPairID := req.Params("election_pair_id")
	if electionPairID == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
		})
	}

	// Get fast election summary with cache
	results, err := api.fastLiveResultUc.GetElectionSummaryWithCache(ctx, electionPairID)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(results), nil
}

// Cache Management Endpoints (Admin Only)

// InvalidateElectionCache godoc
// @Summary Invalidate election cache
// @Description Invalidate all cache entries related to a specific election (admin only)
// @Tags Fast Live Results Admin
// @Accept json
// @Produce json
// @Param election_pair_id path string true "Election Pair ID"
// @Success 200 {object} jsonResponse{data=map[string]interface{}} "Cache invalidated"
// @Router /v1/admin/cache/invalidate/election/{election_pair_id} [post]
func (api *API) InvalidateElectionCache(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.InvalidateElectionCache")
	defer span.End()

	electionPairID := req.Params("election_pair_id")
	if electionPairID == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
		})
	}

	err := api.fastLiveResultUc.InvalidateElectionCache(ctx, electionPairID)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]interface{}{
		"message":     "Election cache invalidated successfully",
		"election_id": electionPairID,
		"timestamp":   time.Now(),
	}), nil
}

// InvalidateCityCache godoc
// @Summary Invalidate city cache
// @Description Invalidate cache entries for a specific city (admin only)
// @Tags Fast Live Results Admin
// @Accept json
// @Produce json
// @Param city_name path string true "City Name"
// @Success 200 {object} jsonResponse{data=map[string]interface{}} "Cache invalidated"
// @Router /v1/admin/cache/invalidate/city/{city_name} [post]
func (api *API) InvalidateCityCache(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.InvalidateCityCache")
	defer span.End()

	cityName := req.Params("city_name")
	if cityName == "" {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "city name is required",
			Code:    400,
		})
	}

	err := api.fastLiveResultUc.InvalidateCityCache(ctx, cityName)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]interface{}{
		"message":   "City cache invalidated successfully",
		"city_name": cityName,
		"timestamp": time.Now(),
	}), nil
}

// InvalidateAllCache godoc
// @Summary Invalidate all cache
// @Description Invalidate all live results cache entries (admin only)
// @Tags Fast Live Results Admin
// @Accept json
// @Produce json
// @Success 200 {object} jsonResponse{data=map[string]interface{}} "All cache invalidated"
// @Router /v1/admin/cache/invalidate/all [post]
func (api *API) InvalidateAllCache(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.InvalidateAllCache")
	defer span.End()

	err := api.fastLiveResultUc.InvalidateAllCache(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]interface{}{
		"message":   "All cache invalidated successfully",
		"timestamp": time.Now(),
	}), nil
}

// GetCacheStatistics godoc
// @Summary Get cache statistics
// @Description Get performance statistics for the live results cache (admin only)
// @Tags Fast Live Results Admin
// @Accept json
// @Produce json
// @Success 200 {object} jsonResponse{data=response.CacheStatisticsResponse} "Cache statistics"
// @Router /v1/admin/cache/statistics [get]
func (api *API) GetCacheStatistics(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.GetCacheStatistics")
	defer span.End()

	stats, err := api.fastLiveResultUc.GetCacheStatistics(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(stats), nil
}

// TriggerLiveResultsBroadcast godoc
// @Summary Trigger live results broadcast
// @Description Manually trigger a broadcast of live results for testing or manual refresh (admin only)
// @Tags Fast Live Results Admin
// @Accept json
// @Produce json
// @Param election_pair_id query string false "Election Pair ID to broadcast"
// @Success 200 {object} jsonResponse{data=map[string]interface{}} "Broadcast triggered"
// @Router /v1/admin/broadcast/live-results [post]
func (api *API) TriggerLiveResultsBroadcast(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultController.TriggerLiveResultsBroadcast")
	defer span.End()

	electionPairID := req.Query("election_pair_id", "")

	connectedClients := api.fastLiveResultUc.GetConnectedClientsCount(ctx)
	if connectedClients == 0 {
		return rest.NewJSONResponse().SetData(map[string]interface{}{
			"message": "No WebSocket clients connected",
			"clients": 0,
		}), nil
	}

	var err error
	if electionPairID != "" {
		err = api.fastLiveResultUc.BroadcastLiveResultsUpdate(ctx, electionPairID)
	} else {
		// Broadcast for all active elections
		allElections, getErr := api.fastLiveResultUc.GetAllElectionsSummaryWithCache(ctx)
		if getErr != nil {
			return custresp.CustomErrorResponse(getErr)
		}

		for _, election := range allElections.Elections {
			if broadcastErr := api.fastLiveResultUc.BroadcastLiveResultsUpdate(ctx, election.ElectionID); broadcastErr != nil {
				// Log error but continue with other elections
			}
		}
	}

	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]interface{}{
		"message":     "Live results broadcast triggered successfully",
		"election_id": electionPairID,
		"clients":     connectedClients,
		"timestamp":   time.Now(),
	}), nil
}
