// internal/infrastructures/websocket/handler.go
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

	welcome := &LiveMessage{
		Type:      MessageTypeHeartbeat,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"status":    "connected",
			"client_id": clientID,
			"message":   "WebSocket connection established",
		},
	}

	if welcomeBytes, err := json.Marshal(welcome); err == nil {
		select {
		case client.Send <- welcomeBytes:
			log.WithFields(log.Fields{
				"client_id": clientID,
			}).Info("[WebSocketHandler] Welcome message sent")
		default:
			log.WithFields(log.Fields{
				"client_id": clientID,
			}).Warn("[WebSocketHandler] Could not send welcome message")
		}
	}

	go h.writePump(client)
	h.readPump(client)
}

func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.unregister <- client
		client.Conn.Close()
		log.WithFields(log.Fields{
			"client_id": client.ID,
		}).Info("[WebSocketHandler] Read pump stopped and client disconnected")
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	client.Conn.SetPongHandler(func(string) error {
		client.UpdateLastSeen()
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		log.WithFields(log.Fields{
			"client_id": client.ID,
		}).Debug("[WebSocketHandler] Pong received, deadline reset")
		return nil
	})

	client.Conn.SetPingHandler(func(message string) error {
		client.UpdateLastSeen()
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Send pong response
		client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		err := client.Conn.WriteMessage(websocket.PongMessage, []byte(message))
		if err != nil {
			log.WithFields(log.Fields{
				"client_id": client.ID,
				"error":     err,
			}).Error("[WebSocketHandler] Failed to send pong")
		} else {
			log.WithFields(log.Fields{
				"client_id": client.ID,
			}).Debug("[WebSocketHandler] Ping received, pong sent")
		}
		return err
	})

	for {
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithFields(log.Fields{
					"client_id": client.ID,
					"error":     err,
				}).Error("[WebSocketHandler] Unexpected close error")
			} else {
				log.WithFields(log.Fields{
					"client_id": client.ID,
					"error":     err,
				}).Debug("[WebSocketHandler] Connection closed")
			}
			break
		}

		client.UpdateLastSeen()
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		log.WithFields(log.Fields{
			"client_id":    client.ID,
			"message_type": messageType,
			"message_size": len(message),
		}).Debug("[WebSocketHandler] Message received")

		if messageType == websocket.TextMessage {
			h.handleTextMessage(client, message)
		} else if messageType == websocket.BinaryMessage {
			// Handle binary messages if needed
			log.WithFields(log.Fields{
				"client_id": client.ID,
			}).Debug("[WebSocketHandler] Binary message received")
		}
	}
}

func (h *Handler) writePump(client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
		log.WithFields(log.Fields{
			"client_id": client.ID,
		}).Info("[WebSocketHandler] Write pump stopped")
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel was closed
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

			log.WithFields(log.Fields{
				"client_id":    client.ID,
				"message_size": len(message),
			}).Debug("[WebSocketHandler] Message sent to client")

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithFields(log.Fields{
					"client_id": client.ID,
					"error":     err,
				}).Debug("[WebSocketHandler] Ping failed, connection likely closed")
				return
			}

			log.WithFields(log.Fields{
				"client_id": client.ID,
			}).Debug("[WebSocketHandler] Ping sent")
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

		errorMsg := &LiveMessage{
			Type:      "error",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"error":   "Invalid message format",
				"details": err.Error(),
			},
		}
		h.hub.sendToClient(client, errorMsg)
		return
	}

	log.WithFields(log.Fields{
		"client_id":    client.ID,
		"message_type": subMsg.Type,
	}).Debug("[WebSocketHandler] Processing subscription message")

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

		errorMsg := &LiveMessage{
			Type:      "error",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"error":        "Unknown message type",
				"message_type": subMsg.Type,
			},
		}
		h.hub.sendToClient(client, errorMsg)
	}
}

func (h *Handler) handleSubscribe(client *Client, subMsg *SubscriptionMessage) {
	filter := &MessageFilter{
		ElectionPairID: subMsg.ElectionPairID,
		City:           subMsg.City,
	}

	client.AddSubscription(subMsg.Subscription, filter)

	log.WithFields(log.Fields{
		"client_id":        client.ID,
		"subscription":     subMsg.Subscription,
		"election_pair_id": subMsg.ElectionPairID,
		"city":             subMsg.City,
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
		HandshakeTimeout:  60 * time.Second, // Match KrakenD timeout
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: false,
		Origins:           []string{"*"},
	})
}

func (h *Handler) WebSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"user_agent": c.Get("User-Agent"),
			"connection": c.Get("Connection"),
			"upgrade":    c.Get("Upgrade"),
			"ws_key":     c.Get("Sec-WebSocket-Key"),
			"ws_version": c.Get("Sec-WebSocket-Version"),
		}).Debug("[WebSocketHandler] WebSocket middleware check")

		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		log.WithFields(log.Fields{
			"method": c.Method(),
			"path":   c.Path(),
		}).Warn("[WebSocketHandler] Non-WebSocket request to WebSocket endpoint")

		return fiber.ErrUpgradeRequired
	}
}
