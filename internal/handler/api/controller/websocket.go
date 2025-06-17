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
	handler      *websocket.Handler
	liveResultUc usecases.LiveResultUseCases
}

type WebSocketControllerOptions struct {
	Handler      *websocket.Handler
	LiveResultUc usecases.LiveResultUseCases
}

func NewWebSocketController(opts *WebSocketControllerOptions) *WebSocketController {
	return &WebSocketController{
		handler:      opts.Handler,
		liveResultUc: opts.LiveResultUc,
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

	connectedClients := wsc.liveResultUc.GetConnectedClientsCount(ctx)

	status := map[string]interface{}{
		"status":             "active",
		"connected_clients":  connectedClients,
		"websocket_endpoint": "/v1/live/ws",
		"supported_subscriptions": []string{
			"all",
			"live_results",
			"city",
			"rankings",
			"summary",
		},
		"message_types": []string{
			"heartbeat",
			"live_results_update",
			"city_results_update",
			"rankings_update",
			"election_summary_update",
			"all_election_update",
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
// @Param city query string false "City name for fast results broadcast"
// @Param type query string false "Broadcast type: vote, election, region, statistics, all" default(all)
// @Success 200 {object} jsonResponse{data=map[string]any} "Broadcast triggered"
// @Router /v1/live/broadcast [post]
func (wsc *WebSocketController) TriggerBroadcast(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "WebSocketController.TriggerBroadcast")
	defer span.End()

	electionPairID := req.Query("election_pair_id", "")
	broadcastType := req.Query("type", "all")
	city := req.Query("city", "")

	connectedClients := wsc.liveResultUc.GetConnectedClientsCount(ctx)
	if connectedClients == 0 {
		return rest.NewJSONResponse().SetData(map[string]interface{}{
			"message": "No connected clients, broadcast skipped",
			"clients": 0,
		}), nil
	}

	var err error
	switch broadcastType {
	case "live_results":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast live results broadcast",
				Code:    400,
			})
		}
		err = wsc.liveResultUc.BroadcastLiveResultsUpdate(ctx, electionPairID)

	case "city":
		if city == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "region is required for fast city results broadcast",
				Code:    400,
			})
		}
		err = wsc.liveResultUc.BroadcastCityResultsUpdate(ctx, city)

	case "rankings":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast rankings broadcast",
				Code:    400,
			})
		}

		err = wsc.liveResultUc.BroadcastRankingsUpdate(ctx, electionPairID, 10)

	case "summary":
		if electionPairID == "" {
			return custresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "election_pair_id is required for fast summary broadcast",
				Code:    400,
			})
		}
		err = wsc.liveResultUc.BroadcastElectionSummaryUpdate(ctx, electionPairID)
	case "all":
		err = wsc.liveResultUc.BroadcastBulkUpdate(ctx, electionPairID)
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
		"city":             city,
		"clients":          connectedClients,
	}), nil
}

func (wsc *WebSocketController) HandleWebSocket(c *fiber.Ctx) error {
	return wsc.handler.UpgradeHandler()(c)
}

func (wsc *WebSocketController) WebSocketMiddleware() fiber.Handler {
	return wsc.handler.WebSocketMiddleware()
}
