package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// server -> client unary communication
type WebSocketConnection struct {
	*websocket.Conn
}

// data received from client
type WsPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	UserName    string              `json:"username"`
	UserID      int                 `json:"user_id"`
	MessageType string              `json:"message_type"`
	Conn        WebSocketConnection `json:"-"`
}

// response sent back to client
type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// parameters to upgrade http connection to websocket
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[WebSocketConnection]string)

var wsChan = make(chan WsPayload)

// websocket handler
func (app *application) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil) // upgrade http connection to websocket
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Printf(fmt.Sprintf("Client connected from %s\n", r.RemoteAddr))

	var response WsJsonResponse
	response.Message = "Connected to server"

	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = "" // register client

	go app.ListenForWS(&conn) // keep listening while websocket connection is alive
}

func (app *application) ListenForWS(conn *WebSocketConnection) {
	// ensure the application recovers gracefully
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("ERROR:", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
			break
		} else {
			payload.Conn = *conn
			wsChan <- payload // push payload from client to channel
		}
	}
}

/*================================================================================*/

func (app *application) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan // pop payload from channel
		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Your account has been deleted"
			response.UserID = e.UserID
			app.broadcastToAll(response) // broadcast response
		default:
		}
	}
}

func (app *application) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		// broadcast to every connected client
		err := client.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("Websocket error on %s: %s\n", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
