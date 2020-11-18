package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)

type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	ID string
	playData []interface{}
	gameInfo []interface{}
}

func (c *Client) readMsg() {
	topLayerMsg := make(map[string]interface{})
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		json.Unmarshal(message, &topLayerMsg) // decode the top layer incoming msg
		if topLayerMsg["msg_type"] == "game_version" {
			if topLayerMsg["msg"] == "practice" {
				log.Println("game_version: practice")
			} else if topLayerMsg["msg"] == "actual" {
				log.Println("game_version: actual")
			}
		} else if topLayerMsg["msg_type"] == "register" {
			if topLayerMsg["msg"] == "not_verified" {
				log.Printf("subject not verified %s\n", topLayerMsg["msg"])
			} else {
				log.Printf("subject %s verified\n", topLayerMsg["msg"])
				c.ID = topLayerMsg["msg"].(string)
			}
		} else if topLayerMsg["msg_type"] == "play_data" {
			c.playData = append(c.playData, topLayerMsg["msg"])
			log.Printf("receive message play_data from participant ID %s\n", c.ID)

		} else if topLayerMsg["msg_type"] == "game_information" {
			c.gameInfo = append(c.gameInfo, topLayerMsg["msg"])
			log.Printf("receive message game_information from participant ID %s\n", c.ID)

		} else {
			log.Fatalf("Warning, unsupported event: %s", topLayerMsg["msg_type"])
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, ID: "default_" + string(rand.Intn(1000))}
	client.hub.register <- client

	go client.readMsg()
}
