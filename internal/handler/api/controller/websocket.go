package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/infrastructures/custresp"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
	"github.com/nocturna-ta/result/internal/usecases"
)

type WebSocketController struct {
	handler           *websocket.Handler
	liveResultService usecases.LiveResultUsecases
	fastResultUc      usecases.FastLiveResultUseCases
}

type WebSocketControllerOptions struct {
	Handler           *websocket.Handler
	LiveResultService usecases.LiveResultUsecases
	FastResultUc      usecases.FastLiveResultUseCases
}

func NewWebSocketController(opts *WebSocketControllerOptions) *WebSocketController {
	return &WebSocketController{
		handler:           opts.Handler,
		liveResultService: opts.LiveResultService,
		fastResultUc:      opts.FastResultUc,
	}
}

// GetLiveResultsStatus godoc
// @Summary Get live results WebSocket status
// @Description Get status information about the live results WebSocket service
// @Tags Live Results
// @Accept json
// @Produce json
// @Success 200 {object} jsonResponse{data=map[string]any} "WebSocket status"
// @Router /v1/live/status [get]
func (wsc *WebSocketController) GetLiveResultsStatus(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "WebSocketController.GetLiveResultsStatus")
	defer span.End()

	connectedClients := wsc.liveResultService.GetConnectedClients(ctx)

	status := map[string]interface{}{
		"status":             "active",
		"connected_clients":  connectedClients,
		"websocket_endpoint": "/v1/live/ws",
		"supported_subscriptions": []string{
			"all",
			"election",
			"region",
			"statistics",
			"fast_live_results",
			"fast_city",
			"fast_rankings",
			"fast_summary",
		},
		"message_types": []string{
			"vote_update",
			"election_update",
			"region_update",
			"statistics_update",
			"heartbeat",
			"fast_live_results_update",
			"fast_city_results_update",
			"fast_rankings_update",
			"fast_election_summary_update",
			"fast_all_election_update",
		},
	}

	return rest.NewJSONResponse().SetData(status), nil
}

// TriggerBroadcast godoc
// @Summary Trigger manual broadcast
// @Description Manually trigger a broadcast of current results (for testing/admin purposes)
// @Tags Live Results
// @Accept json
// @Produce json
// @Param election_pair_id query string false "Election Pair ID to broadcast"
// @Param region query string false "Region to broadcast"
// @Param type query string false "Broadcast type: vote, election, region, statistics, all" default(all)
// @Success 200 {object} jsonResponse{data=map[string]any} "Broadcast triggered"
// @Router /v1/live/broadcast [post]
func (wsc *WebSocketController) TriggerBroadcast(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "WebSocketController.TriggerBroadcast")
	defer span.End()

	electionPairID := req.Query("election_pair_id", "")
	region := req.Query("region", "")
	broadcastType := req.Query("type", "all")

	connectedClients := wsc.liveResultService.GetConnectedClients(ctx)
	if connectedClients == 0 {
		return rest.NewJSONResponse().SetData(map[string]interface{}{
			"message": "No connected clients, broadcast skipped",
			"clients": 0,
		}), nil
	}

	var err error
	switch broadcastType {
	case "statistics":
		err = wsc.liveResultService.BroadcastStatisticsUpdate(ctx)
	case "election":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for election broadcast",
				Code:    400,
			})
		}
		err = wsc.liveResultService.BroadcastElectionUpdate(ctx, electionPairID)
	case "region":
		if region == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "region is required for region broadcast",
				Code:    400,
			})
		}
		err = wsc.liveResultService.BroadcastRegionUpdate(ctx, region)
	case "fast_live_results":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast live results broadcast",
				Code:    400,
			})
		}
		err = wsc.fastResultUc.BroadcastLiveResultsUpdate(ctx, electionPairID)

	case "fast_city":
		if region == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "region is required for fast city results broadcast",
				Code:    400,
			})
		}
		err = wsc.fastResultUc.BroadcastCityResultsUpdate(ctx, region)

	case "fast_rankings":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast rankings broadcast",
				Code:    400,
			})
		}

		err = wsc.fastResultUc.BroadcastRankingsUpdate(ctx, electionPairID, 10)

	case "fast_summary":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast summary broadcast",
				Code:    400,
			})
		}
		err = wsc.fastResultUc.BroadcastElectionSummaryUpdate(ctx, electionPairID)
	case "all":
		err = wsc.liveResultService.BroadcastAllUpdates(ctx, electionPairID, region)
	default:
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid broadcast type",
			Code:    400,
		})
	}

	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(map[string]interface{}{
		"message":          "Broadcast triggered successfully",
		"type":             broadcastType,
		"election_pair_id": electionPairID,
		"region":           region,
		"clients":          connectedClients,
	}), nil
}

func (wsc *WebSocketController) HandleWebSocket(c *fiber.Ctx) error {
	return wsc.handler.UpgradeHandler()(c)
}

func (wsc *WebSocketController) WebSocketMiddleware() fiber.Handler {
	return wsc.handler.WebSocketMiddleware()
}
