package websocket

import (
	"context"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"strings"
	"sync"
	"time"
)

type MessageType string

const (
	MessageTypeVoteUpdate  MessageType = "vote_update"
	MessageTypeHeartbeat   MessageType = "heartbeat"
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"

	MessageTypeLiveResults         MessageType = "live_results_update"
	MessageTypeIncremental         MessageType = "incremental_update"
	MessageTypeLiveCityResults     MessageType = "city_results_update"
	MessageTypeLiveRankings        MessageType = "rankings_update"
	MessageTypeLiveElectionSummary MessageType = "election_summary_update"
	MessageTypeLiveAllElections    MessageType = "all_elections_update"
)

type SubscriptionType string

const (
	SubscriptionAll         SubscriptionType = "all"
	SubscriptionLiveResults SubscriptionType = "live_results"
	SubscriptionCity        SubscriptionType = "city"
	SubscriptionRankings    SubscriptionType = "rankings"
	SubscriptionSummary     SubscriptionType = "summary"
)

type LiveMessage struct {
	Type      MessageType    `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      interface{}    `json:"data,omitempty"`
	Filter    *MessageFilter `json:"filter,omitempty"`
}

type MessageFilter struct {
	ElectionPairID string `json:"election_pair_id,omitempty"`
	City           string `json:"city,omitempty"`
}

type SubscriptionMessage struct {
	Type           MessageType      `json:"type"`
	Subscription   SubscriptionType `json:"subscription"`
	ElectionPairID string           `json:"election_pair_id,omitempty"`
	City           string           `json:"city,omitempty"`
}

type Client struct {
	ID            string
	Conn          *websocket.Conn
	Send          chan []byte
	Subscriptions map[SubscriptionType]*MessageFilter
	LastSeen      time.Time
	mu            sync.RWMutex
}

type Hub struct {
	clients    map[string]*Client
	broadcast  chan *LiveMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewHub(ctx context.Context) *Hub {
	hubCtx, cancel := context.WithCancel(ctx)
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *LiveMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		ctx:        hubCtx,
		cancel:     cancel,
	}
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:            id,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		Subscriptions: make(map[SubscriptionType]*MessageFilter),
		LastSeen:      time.Now(),
	}
}

func (c *Client) UpdateLastSeen() {
	c.mu.Lock()
	c.LastSeen = time.Now()
	c.mu.Unlock()
}

func (c *Client) AddSubscription(subType SubscriptionType, filter *MessageFilter) {
	c.mu.Lock()
	c.Subscriptions[subType] = filter
	c.mu.Unlock()
}

func (c *Client) RemoveSubscription(subType SubscriptionType) {
	c.mu.Lock()
	delete(c.Subscriptions, subType)
	c.mu.Unlock()
}

func (c *Client) GetSubscriptions() map[SubscriptionType]*MessageFilter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	subs := make(map[SubscriptionType]*MessageFilter)
	for k, v := range c.Subscriptions {
		subs[k] = v
	}
	return subs
}

func (h *Hub) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

			log.WithFields(log.Fields{
				"client_id":     client.ID,
				"total_clients": len(h.clients),
			}).Info("[Websocket Hub] Client registered")

			welcome := &LiveMessage{
				Type:      MessageTypeHeartbeat,
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"status":    "connected",
					"client_id": client.ID,
				},
			}
			h.sendToClient(client, welcome)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

			log.WithFields(log.Fields{
				"client_id":     client.ID,
				"total_clients": len(h.clients),
			}).Info("[Websocket Hub] Client unregistered")

		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				if h.shouldSendToClient(client, msg) {
					select {
					case client.Send <- h.messageToBytes(msg):
					default:
						h.mu.RUnlock()
						h.unregister <- client
						h.mu.RLock()
					}
				}
			}
			h.mu.RUnlock()
		case <-ticker.C:
			h.sendHeartbeat()
			h.cleanupStaleConnections()
		}
	}
}

func (h *Hub) shouldSendToClient(client *Client, message *LiveMessage) bool {
	client.mu.RLock()
	defer client.mu.RUnlock()

	if message.Type == MessageTypeHeartbeat {
		return true
	}

	if len(client.Subscriptions) == 0 {
		return false
	}

	for subType, filter := range client.Subscriptions {
		if h.matchesSubscription(subType, filter, message, client.ID) {
			return true
		}
	}
	return false
}

func (h *Hub) sendToClient(client *Client, message *LiveMessage) {
	select {
	case client.Send <- h.messageToBytes(message):
	default:
		h.unregister <- client
	}
}

func (h *Hub) messageToBytes(message *LiveMessage) []byte {
	data, _ := json.Marshal(message)
	return data
}

func (h *Hub) sendHeartbeat() {
	heartbeat := &LiveMessage{
		Type:      MessageTypeHeartbeat,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"status":  "alive",
			"clients": len(h.clients),
		},
	}

	h.mu.RLock()
	for _, client := range h.clients {
		h.sendToClient(client, heartbeat)
	}
	h.mu.RUnlock()
}

func (h *Hub) cleanupStaleConnections() {
	now := time.Now()
	staleThreshold := 2 * time.Minute

	h.mu.RLock()
	var staleClients []*Client
	for _, client := range h.clients {
		if now.Sub(client.LastSeen) > staleThreshold {
			staleClients = append(staleClients, client)
		}
	}
	h.mu.RUnlock()

	for _, client := range staleClients {
		h.unregister <- client
	}
}

func (h *Hub) BroadcastLiveResults(results *response.ElectionResultsResponse) {
	message := &LiveMessage{
		Type:      MessageTypeLiveResults,
		Timestamp: time.Now(),
		Data:      results,
		Filter: &MessageFilter{
			ElectionPairID: results.ElectionID,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping fast live results update message")
	}
}

func (h *Hub) BroadcastCityResults(results *response.CityResultsResponse) {
	message := &LiveMessage{
		Type:      MessageTypeLiveCityResults,
		Timestamp: time.Now(),
		Data:      results,
		Filter: &MessageFilter{
			City: results.CityName,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping fast city results update message")
	}
}

func (h *Hub) BroadcastRankings(results *response.CityRankingsResponse) {
	message := &LiveMessage{
		Type:      MessageTypeLiveRankings,
		Timestamp: time.Now(),
		Data:      results,
		Filter: &MessageFilter{
			ElectionPairID: results.ElectionID,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping fast rankings update message")
	}
}

func (h *Hub) BroadcastElectionSummary(results *response.ElectionSummaryResponse) {
	message := &LiveMessage{
		Type:      MessageTypeLiveElectionSummary,
		Timestamp: time.Now(),
		Data:      results,
		Filter: &MessageFilter{
			ElectionPairID: results.ElectionID,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping fast election summary update message")
	}
}

func (h *Hub) BroadcastFastAllElections(results *response.AllElectionsSummaryResponse) {
	message := &LiveMessage{
		Type:      MessageTypeLiveAllElections,
		Timestamp: time.Now(),
		Data:      results,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping fast all elections update message")
	}
}

func (h *Hub) BroadcastVoteUpdate(results *response.VoteResultResponse) {
	message := &LiveMessage{
		Type:      MessageTypeVoteUpdate,
		Timestamp: time.Now(),
		Data:      results,
		Filter: &MessageFilter{
			ElectionPairID: results.ElectionPairID,
			City:           results.Region,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping vote update message")
	}

}

func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) Stop() {
	h.cancel()

	h.mu.Lock()
	for _, client := range h.clients {
		close(client.Send)
	}
	h.clients = make(map[string]*Client)
	h.mu.Unlock()
}

func (h *Hub) matchesSubscription(subType SubscriptionType, filter *MessageFilter, message *LiveMessage, clientID string) bool {
	switch subType {
	case SubscriptionAll:
		log.WithFields(log.Fields{
			"client_id":    clientID,
			"subscription": "all",
		}).Debug("[WebSocket] Sending to 'all' subscription")
		return true

	case SubscriptionLiveResults:
		return h.matchesFastLiveResults(filter, message, clientID)

	case SubscriptionCity:
		return h.matchesFastCity(filter, message, clientID)

	case SubscriptionRankings:
		return h.matchesFastRankings(filter, message, clientID)

	case SubscriptionSummary:
		return h.matchesFastSummary(filter, message, clientID)
	}

	return false
}

func (h *Hub) matchesFastLiveResults(filter *MessageFilter, message *LiveMessage, clientID string) bool {
	if message.Type != MessageTypeLiveResults {
		return false
	}

	if filter == nil || strings.TrimSpace(filter.ElectionPairID) == "" {
		return true
	}

	if message.Filter == nil {
		return false
	}

	return strings.TrimSpace(filter.ElectionPairID) == strings.TrimSpace(message.Filter.ElectionPairID)
}

func (h *Hub) matchesFastCity(filter *MessageFilter, message *LiveMessage, clientID string) bool {
	if message.Type != MessageTypeLiveCityResults {
		return false
	}

	if filter == nil || strings.TrimSpace(filter.City) == "" {
		return true
	}

	if message.Filter == nil {
		return false
	}

	clientRegion := strings.TrimSpace(strings.ToLower(filter.City))
	messageRegion := strings.TrimSpace(strings.ToLower(message.Filter.City))

	return clientRegion == messageRegion
}

func (h *Hub) matchesFastRankings(filter *MessageFilter, message *LiveMessage, clientID string) bool {
	if message.Type != MessageTypeLiveRankings {
		return false
	}

	if filter == nil || strings.TrimSpace(filter.ElectionPairID) == "" {
		return true
	}

	if message.Filter == nil {
		return false
	}

	return strings.TrimSpace(filter.ElectionPairID) == strings.TrimSpace(message.Filter.ElectionPairID)
}

func (h *Hub) matchesFastSummary(filter *MessageFilter, message *LiveMessage, clientID string) bool {
	if message.Type != MessageTypeLiveElectionSummary && message.Type != MessageTypeLiveAllElections {
		return false
	}

	if filter == nil || strings.TrimSpace(filter.ElectionPairID) == "" {
		return true
	}

	if message.Filter == nil {
		return false
	}

	return strings.TrimSpace(filter.ElectionPairID) == strings.TrimSpace(message.Filter.ElectionPairID)
}
