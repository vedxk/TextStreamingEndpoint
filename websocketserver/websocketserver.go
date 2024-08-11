package websocketserver

import (
	"errors"
	"log"
	"myapp/inference"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = make(map[*websocket.Conn]bool)
	connMu      sync.Mutex
)

const (
	maxRetries      = 3
	responseTimeout = 5 * time.Second
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	connMu.Lock()
	connections[conn] = true
	connMu.Unlock()

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer func() {
		connMu.Lock()
		delete(connections, conn)
		connMu.Unlock()

		if err := conn.Close(); err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			log.Println("Connection close error:", err)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Ping error:", err)
				return
			}
		default:
			conn.SetReadDeadline(time.Now().Add(pongWait)) // Reset deadline for each read
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Check if the error is because the connection is closing normally
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("Connection closing normally:", err)
					return
				}
				log.Println("Read error:", err)
				// No need to send an error message if the connection is closing or closed
				if conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)) != nil {
					return
				}
				sendErrorResponse(conn, "Error reading message from WebSocket.")
				return
			}

			response, err := getResponseWithRetry(string(message), maxRetries)
			if err != nil {
				log.Println("Inference error:", err)
				if conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)) != nil {
					return
				}
				sendErrorResponse(conn, err.Error())
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("Connection closing normally:", err)
					return
				}
				log.Println("Write error:", err)
				if conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)) != nil {
					return
				}
				sendErrorResponse(conn, "Error sending message to WebSocket.")
				return
			}
		}
	}
}

func getResponseWithRetry(prompt string, retries int) (string, error) {
	for i := 0; i < retries; i++ {
		responseChan := make(chan string)
		errorChan := make(chan error)

		go func() {
			response, err := inference.GetResponse(prompt)
			if err != nil {
				errorChan <- err
			} else {
				responseChan <- response
			}
		}()

		select {
		case response := <-responseChan:
			return response, nil
		case err := <-errorChan:
			log.Println("Provider error:", err)
			if i == retries-1 {
				return "", errors.New("all retries failed: " + err.Error())
			}
		case <-time.After(responseTimeout):
			log.Println("Response timeout")
			if i == retries-1 {
				return "", errors.New("response timeout: all retries failed")
			}
		}
	}

	return "", errors.New("unknown error")
}

func sendErrorResponse(conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("Failed to send error message:", err)
		}
	}
}
