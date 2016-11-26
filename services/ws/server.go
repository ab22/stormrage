package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketServer interface {
	OnConnect(http.ResponseWriter, *http.Request) error
	AddClient(WebsocketClient)
	RemoveClient(WebsocketClient)
	LogError(error)
}

// Server contains all information to host the websocket server.
type websocketServer struct {
	messages       []string
	clients        map[int]WebsocketClient
	addClientCh    chan WebsocketClient
	removeClientCh chan WebsocketClient
	errorCh        chan error
	upgrader       websocket.Upgrader
}

// NewServer initializes a new Client struct.
func NewServer() WebsocketServer {
	server := &websocketServer{
		messages:       []string{},
		clients:        make(map[int]WebsocketClient),
		addClientCh:    make(chan WebsocketClient),
		removeClientCh: make(chan WebsocketClient),
		errorCh:        make(chan error),

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	go server.Listen()

	return server
}

// AddClient adds a new client to the server's client list. AddClient is a
// blocking function.
func (s *websocketServer) AddClient(client WebsocketClient) {
	s.addClientCh <- client
}

// RemoveClient removes an existing client from the server's client list.
// RemoveClient is a blocking function.
func (s *websocketServer) RemoveClient(client WebsocketClient) {
	s.removeClientCh <- client
}

// LogError lets other goroutines to log errors though the server. In the
// future, error processing might be done in this function so it is better
// to call this function when an error happened. This function is a blocking
// function.
func (s *websocketServer) LogError(err error) {
	s.errorCh <- err
}

// OnConnect 'upgrades' a normal HTTP request to a websocket connection.
func (s *websocketServer) OnConnect(w http.ResponseWriter, r *http.Request) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	client := NewClient(conn, s)
	s.AddClient(client)
	client.Listen()
	s.RemoveClient(client)

	return nil
}

// Listen is the main server loop. This function awaits for incoming messages
// from all channels.
func (s *websocketServer) Listen() {
	for {
		select {
		case client := <-s.addClientCh:
			s.addClient(client)

		case client := <-s.removeClientCh:
			s.removeClient(client)

		case err := <-s.errorCh:
			log.Println("server error channel:", err.Error())
		}
	}
}

func (s *websocketServer) addClient(client WebsocketClient) {
	s.clients[client.GetID()] = client
	log.Println("Client has joined the channel!")
}

func (s *websocketServer) removeClient(client WebsocketClient) {
	delete(s.clients, client.GetID())
	log.Println("Client has left the channel!")
}

func (s *websocketServer) broadcastMessage(message []byte) {
	for _, client := range s.clients {
		client.Write(message)
	}
}
