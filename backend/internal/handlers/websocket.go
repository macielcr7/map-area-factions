package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketHandler struct {
	clients map[*websocket.Conn]bool
	hub     chan []byte
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type FactionUpdateMessage struct {
	FactionID string `json:"faction_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Action    string `json:"action"` // "created", "updated", "deleted"
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients: make(map[*websocket.Conn]bool),
		hub:     make(chan []byte),
	}
}

func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
	defer func() {
		delete(h.clients, c)
		c.Close()
	}()

	// Register new client
	h.clients[c] = true

	// Send welcome message
	welcomeMsg := WSMessage{
		Type: "connected",
		Payload: map[string]string{
			"message": "Connected to Map Factions real-time updates",
		},
	}
	
	if data, err := json.Marshal(welcomeMsg); err == nil {
		c.WriteMessage(websocket.TextMessage, data)
	}

	// Listen for messages
	for {
		var msg WSMessage
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}
		
		// Handle different message types
		switch msg.Type {
		case "ping":
			pongMsg := WSMessage{
				Type: "pong",
				Payload: map[string]string{
					"message": "pong",
				},
			}
			c.WriteJSON(pongMsg)
		case "subscribe":
			// Handle subscription to specific channels
			log.Printf("Client subscribed to: %v", msg.Payload)
		}
	}
}

func (h *WebSocketHandler) BroadcastFactionUpdate(factionID, name, color, action string) {
	message := WSMessage{
		Type: "faction_update",
		Payload: FactionUpdateMessage{
			FactionID: factionID,
			Name:      name,
			Color:     color,
			Action:    action,
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling faction update: %v", err)
		return
	}

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			delete(h.clients, client)
			client.Close()
		}
	}
}

func (h *WebSocketHandler) BroadcastGeometryUpdate(geometryID, factionID, action string) {
	message := WSMessage{
		Type: "geometry_update",
		Payload: map[string]string{
			"geometry_id": geometryID,
			"faction_id":  factionID,
			"action":      action,
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling geometry update: %v", err)
		return
	}

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			delete(h.clients, client)
			client.Close()
		}
	}
}

// HTTP upgrade endpoint
func (h *WebSocketHandler) UpgradeHandler(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}