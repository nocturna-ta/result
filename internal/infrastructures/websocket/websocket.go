package websocket

import (
	"context"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"sync"
	"time"
)

type MessageType string

const (
	MessageTypeVoteUpdate     MessageType = "vote_update"
	MessageTypeElectionUpdate MessageType = "election_update"
	MessageTypeRegionUpdate   MessageType = "region_update"
	MessageTypeStatistics     MessageType = "statistics_update"
	MessageTypeHeartbeat      MessageType = "heartbeat"
	MessageTypeSubscribe      MessageType = "subscribe"
	MessageTypeUnsubscribe    MessageType = "unsubscribe"
)

type SubscriptionType string

const (
	SubscriptionAll        SubscriptionType = "all"
	SubscriptionElection   SubscriptionType = "election"
	SubscriptionRegion     SubscriptionType = "region"
	SubscriptionStatistics SubscriptionType = "statistics"
)

type LiveMessage struct {
	Type      MessageType    `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      interface{}    `json:"data,omitempty"`
	Filter    *MessageFilter `json:"filter,omitempty"`
}

type MessageFilter struct {
	ElectionPairID string `json:"election_pair_id,omitempty"`
	Region         string `json:"region,omitempty"`
}

type SubscriptionMessage struct {
	Type           MessageType      `json:"type"`
	Subscription   SubscriptionType `json:"subscription"`
	ElectionPairID string           `json:"election_pair_id,omitempty"`
	Region         string           `json:"region,omitempty"`
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
		switch subType {
		case SubscriptionAll:
			return true

		case SubscriptionElection:
			if message.Type == MessageTypeElectionUpdate || message.Type == MessageTypeVoteUpdate {
				if filter == nil || filter.ElectionPairID == "" {
					return true
				}
				if message.Filter != nil && message.Filter.ElectionPairID == filter.ElectionPairID {
					return true
				}
			}

		case SubscriptionRegion:
			if message.Type == MessageTypeRegionUpdate || message.Type == MessageTypeVoteUpdate {
				if filter == nil || filter.Region == "" {
					return true
				}
				if message.Filter != nil && message.Filter.Region == filter.Region {
					return true
				}
			}
		case SubscriptionStatistics:
			if message.Type == MessageTypeStatistics {
				return true
			}
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

func (h *Hub) BroadcastVoteUpdate(voteResult *response.VoteResultResponse) {
	message := &LiveMessage{
		Type:      MessageTypeVoteUpdate,
		Timestamp: time.Now(),
		Data:      voteResult,
		Filter: &MessageFilter{
			ElectionPairID: voteResult.ElectionPairID,
			Region:         voteResult.Region,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping vote update message")
	}
}

func (h *Hub) BroadcastElectionUpdate(electionResult *response.ElectionVoteResultResponse) {
	message := &LiveMessage{
		Type:      MessageTypeElectionUpdate,
		Timestamp: time.Now(),
		Data:      electionResult,
		Filter: &MessageFilter{
			ElectionPairID: electionResult.ElectionPairID,
			Region:         electionResult.Region,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping election update message")
	}
}

func (h *Hub) BroadcastRegionUpdate(regionResult *response.RegionVoteResultResponse) {
	message := &LiveMessage{
		Type:      MessageTypeRegionUpdate,
		Timestamp: time.Now(),
		Data:      regionResult,
		Filter: &MessageFilter{
			Region: regionResult.Region,
		},
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping region update message")
	}
}

func (h *Hub) BroadcastStatisticsUpdate(stats *response.VoteStatisticsResponse) {
	message := &LiveMessage{
		Type:      MessageTypeStatistics,
		Timestamp: time.Now(),
		Data:      stats,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Warn("[WebSocketHub] Broadcast channel full, dropping statistics update message")
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
