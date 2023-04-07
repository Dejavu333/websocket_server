package websockets

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// ============================================================================================================================
// IWebSocketServer interface
// ============================================================================================================================
type IWebSocketServer interface {
	Start()
	Stop()
	Broadcast(p_channelName string, p_data interface{})
	AddChannel(p_channelName string)
}

// ----------------------------------------------
// IOC
// ----------------------------------------------
func NewIWebSocketServer() IWebSocketServer {
	return NewDefaultWebSocketServer()
}

// ============================================================================================================================
// DefaultWebSocketServer struct implements IWebSocketServer interface
// ============================================================================================================================
type DefaultWebSocketServer struct {
	upgrader *websocket.Upgrader
	clients  map[*websocket.Conn]bool
	channels map[string]map[*websocket.Conn]bool
}

// ----------------------------------------------
// constructors
// ----------------------------------------------
func NewDefaultWebSocketServer() *DefaultWebSocketServer {
	return &DefaultWebSocketServer{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:  make(map[*websocket.Conn]bool),
		channels: make(map[string]map[*websocket.Conn]bool),
	}
}

// ----------------------------------------------
// methods
// ----------------------------------------------
/* starts the server */
func (thisServer *DefaultWebSocketServer) Start() {
	port := utils.getEnvOrDefault("PORT", "8080")
	host := utils.getEnvOrDefault("HOST", "0.0.0.0")
	http.HandleFunc("/", thisServer.handleHTTPRequest)
	logrus.Info("Started websocket server on port 8080...")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		logrus.Fatal("ListenAndServe: ", err)
	}
}

/* stops the server */
func (thisServer *DefaultWebSocketServer) Stop() {

	for client := range thisServer.clients {
		client.Close()
		delete(thisServer.clients, client)
	}
}

/* handles HTTP requests */
func (thisServer *DefaultWebSocketServer) handleHTTPRequest(p_responseWriter http.ResponseWriter, p_request *http.Request) {

	/* checks the request URL and upgrades the connection accordingly */
	if thisServer.channels[p_request.URL.Path] != nil {
		channelName := p_request.URL.Path
		thisServer.upgradeWebSocket(p_responseWriter, p_request, channelName)
	} else {
		http.NotFound(p_responseWriter, p_request)
	}
}

/* upgrades the HTTP connection to a websocket connection */
func (thisServer *DefaultWebSocketServer) upgradeWebSocket(p_responseWriter http.ResponseWriter, p_request *http.Request, p_channel string) {

	conn, err := thisServer.upgrader.Upgrade(p_responseWriter, p_request, nil)
	if err != nil {
		logrus.Println(err)
		return
	}

	/* adds the client to the clients map and the channel map */
	thisServer.clients[conn] = true
	if _, ok := thisServer.channels[p_channel]; !ok {
		thisServer.channels[p_channel] = make(map[*websocket.Conn]bool)
	}
	thisServer.channels[p_channel][conn] = true

	logrus.Println("Client connected:", conn.RemoteAddr())

	/* we could use a goroutine here to listen for messages from the client, or send messages to the client */
}

/* broadcasts the data to all clients in the specified channel */
func (thisServer *DefaultWebSocketServer) Broadcast(p_channel string, p_data interface{}) {

	if clients, ok := thisServer.channels[p_channel]; ok {
		for client := range clients {
			err := client.WriteJSON(p_data)
			if err != nil {
				logrus.Printf("error: %v", err)
				client.Close()
				delete(thisServer.clients, client)
				delete(thisServer.channels[p_channel], client)
			}
		}
	}
}

/* adds a channel to the server */
func (thisServer *DefaultWebSocketServer) AddChannel(p_channelName string) {

	if _, ok := thisServer.channels[p_channelName]; !ok {
		thisServer.channels[p_channelName] = make(map[*websocket.Conn]bool)
	}
}

// ============================================================================================================================
// utility functions
// ============================================================================================================================
func getEnvOrDefault(p_key string, p_defaultValue string) string {
	value := os.Getenv(p_key)
	if value == "" {
		value = p_defaultValue
	}
	return value
}
