package websocket

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/log"
	"time"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) HandleConnection(c *websocket.Conn) {
	clientID := uuid.New().String()
	client := NewClient(clientID, c)

	h.hub.register <- client

	go h.writePump(client)
	h.readPump(client)
}

func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.UpdateLastSeen()
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithFields(log.Fields{
					"client_id": client.ID,
					"error":     err,
				}).Error("[WebSocketHandler] Unexpected close error")
			}
			break
		}

		client.UpdateLastSeen()

		if messageType == websocket.TextMessage {
			h.handleTextMessage(client, message)
		}
	}
}

func (h *Handler) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.WithFields(log.Fields{
					"client_id": client.ID,
					"error":     err,
				}).Error("[WebSocketHandler] Write message error")
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Handler) handleTextMessage(client *Client, message []byte) {
	var subMsg SubscriptionMessage
	if err := json.Unmarshal(message, &subMsg); err != nil {
		log.WithFields(log.Fields{
			"client_id": client.ID,
			"error":     err,
			"message":   string(message),
		}).Error("[WebSocketHandler] Failed to unmarshal subscription message")
		return
	}

	switch subMsg.Type {
	case MessageTypeSubscribe:
		h.handleSubscribe(client, &subMsg)
	case MessageTypeUnsubscribe:
		h.handleUnsubscribe(client, &subMsg)
	default:
		log.WithFields(log.Fields{
			"client_id":    client.ID,
			"message_type": subMsg.Type,
		}).Warn("[WebSocketHandler] Unknown message type")
	}
}

func (h *Handler) handleSubscribe(client *Client, subMsg *SubscriptionMessage) {
	filter := &MessageFilter{
		ElectionPairID: subMsg.ElectionPairID,
		Region:         subMsg.Region,
	}

	client.AddSubscription(subMsg.Subscription, filter)

	log.WithFields(log.Fields{
		"client_id":        client.ID,
		"subscription":     subMsg.Subscription,
		"election_pair_id": subMsg.ElectionPairID,
		"region":           subMsg.Region,
	}).Info("[WebSocketHandler] Client subscribed")

	ack := &LiveMessage{
		Type:      MessageTypeSubscribe,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"subscription": subMsg.Subscription,
			"status":       "subscribed",
			"filter":       filter,
		},
	}

	h.hub.sendToClient(client, ack)
}

func (h *Handler) handleUnsubscribe(client *Client, subMsg *SubscriptionMessage) {
	client.RemoveSubscription(subMsg.Subscription)

	log.WithFields(log.Fields{
		"client_id":    client.ID,
		"subscription": subMsg.Subscription,
	}).Info("[WebSocketHandler] Client unsubscribed")

	ack := &LiveMessage{
		Type:      MessageTypeUnsubscribe,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"subscription": subMsg.Subscription,
			"status":       "unsubscribed",
		},
	}

	h.hub.sendToClient(client, ack)
}

func (h *Handler) UpgradeHandler() fiber.Handler {
	return websocket.New(h.HandleConnection, websocket.Config{
		HandshakeTimeout:  10 * time.Second,
		EnableCompression: true,
	})
}

func (h *Handler) WebSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
