package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/ab22/stormrage/handlers/httputils"
	"github.com/gorilla/websocket"
)

var (
	// Client ID
	clientID = 0

	// Mutex to inc client ID
	clientIDMutex = sync.Mutex{}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 54 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Client's buffered channel size.
	messageChannelSize = 16
)

type WebsocketClient interface {
	GetID() int
	CloseAndWait()
	CloseOrIgnore()
	Listen()
	Write([]byte)
	WriteMany([]string) bool
	WriteAndWait([]byte) bool
	LogError(error)
}

// Client contains all information associated with a websocket client conn.
type websocketClient struct {
	ID         int
	conn       *websocket.Conn
	server     WebsocketServer
	msgCh      chan []byte
	closeCh    chan bool
	pingWriter *pingWriter
}

// generateClientID increments the global client id in a thread safe way.
func generateClientID() int {
	clientIDMutex.Lock()
	newID := clientID
	clientID++
	clientIDMutex.Unlock()

	return newID
}

// NewClient initializes a new Client struct, sets the default read limits and
// deadlines and creates a Pong Handler for the connection.
func NewClient(conn *websocket.Conn, server WebsocketServer) WebsocketClient {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	return &websocketClient{
		ID:         generateClientID(),
		conn:       conn,
		server:     server,
		msgCh:      make(chan []byte, messageChannelSize),
		closeCh:    make(chan bool, 1),
		pingWriter: nil,
	}
}

func (c *websocketClient) GetID() int {
	return c.ID
}

// CloseAndWait sends a signal to the close channel to terminate the
// Read() and Write() client's goroutines. It waits for the message
// to go through before returning.
func (c *websocketClient) CloseAndWait() {
	c.closeCh <- true
}

// CloseOrIgnore sends a signal to the close channel to terminate the
// Read() and Write() client's goroutines. It attempts to send the signal
// to the channel or ignore it right away if it can't go through.
func (c *websocketClient) CloseOrIgnore() {
	select {
	case c.closeCh <- true:
		// Close signal sent successfully.
	default:
		// If close signal can't be sent, it means there's already
		// one signal in queue. Ignore send.
	}
}

// Listen will spawn a new goroutine to write to the client and will make a
// call to onRead() for when the client sends in data.
func (c *websocketClient) Listen() {
	go c.onWrite()
	c.onRead()
}

// LogError is a helper function that calls the internal server pointer to log the error.
func (c *websocketClient) LogError(err error) {
	c.server.LogError(err)
}

// Write attempts to send a message to the client. If the client
// channel's buffer is full, then we spawn a new goroutine with a call
// to client.writeAndWait which will wait waitTime duration. This is done in
// order to prevent blocking the caller for too long.
func (c *websocketClient) Write(msg []byte) {
	select {
	case c.msgCh <- msg:
		// Message sent.

	default:
		// Client's queue is full, so spawn a new goroutine to avoid
		// blocking the caller for long.
		go c.WriteAndWait(msg)

	}
}

// WriteAndWait attempts to send the message to the client's channel. If it
// takes too long to respond, we disconnect the client.
func (c *websocketClient) WriteAndWait(msg []byte) bool {
	select {
	case c.msgCh <- msg:
		// Message successfully sent.
		return true

	case <-time.After(writeWait):
		// If client queue is full and timed out, then send a close channel
		// signal.
		c.CloseOrIgnore()

		return false
	}
}

// WriteMany takes an array of strings and calls the Write() function
// to send each message in sequence to the client's channel. If one of them
// fails, we return.
func (c *websocketClient) WriteMany(msgs []string) bool {
	var ok bool

	for _, msg := range msgs {
		ok = c.WriteAndWait([]byte(msg))

		if !ok {
			return false
		}
	}

	return true
}

// write is a wrapper around gorilla's Connection.WriteMessage function
// that sets a write deadline to timeout a write call and the actual
// WriteMessage call to send the message.
func (c *websocketClient) write(messageType int, msg []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, msg)
}

// onWrite function writes to the client all messages. Messages include server
// custom messages, ping messages and close connection messages. Must be run
// on a different goroutine from onRead().
func (c *websocketClient) onWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.stopPing()
		ticker.Stop()
		c.conn.Close()

		// End onRead() goroutine.
		c.CloseOrIgnore()
	}()

	for {
		select {
		case msg := <-c.msgCh:
			err := c.write(websocket.TextMessage, msg)

			if err != nil {
				c.LogError(err)
				return
			}

		case <-ticker.C:
			err := c.write(websocket.PingMessage, []byte{})

			if err != nil {
				err = fmt.Errorf("client[%d] ping timeout: %s", c.ID, err.Error())
				c.LogError(err)
				return
			}

		case <-c.closeCh:
			c.write(websocket.CloseMessage, []byte{})
			return
		}
	}
}

// onRead function reads all messages from client and makes sure to send them
// to the server. Client messages include the pong messages and custom client
// messages. Must be run on a different goroutine frmo onWrite().
func (c *websocketClient) onRead() {
	defer func() {
		c.conn.Close()

		// End onWrite() goroutine.
		c.CloseOrIgnore()
	}()

	for {
		select {
		case <-c.closeCh:
			return

		default:
			_, reader, err := c.conn.NextReader()

			if err != nil {
				return
			}

			req := &request{}
			err = httputils.DecodeJSON(reader, req)
			if err != nil {
				err = fmt.Errorf("ws client: error decoding json: %v", err)
				c.LogError(err)
				continue
			}

			c.processRequest(req)
		}
	}
}

func (c *websocketClient) startPing(req *request) {
	if !req.IsValidIP() {
		c.Write([]byte("{ \"error\": \"Invalid IP!\"}"))
		return
	}

	if c.pingWriter != nil {
		err := c.pingWriter.Kill()
		c.pingWriter = nil

		if err != nil {
			err = fmt.Errorf("start ping: error killing ping command: %v", err)
			c.LogError(err)
		}

	}

	c.pingWriter = newPingWriter(c, req.IP)
	go c.pingWriter.StartAndWait()
}

func (c *websocketClient) stopPing() {
	if c.pingWriter == nil {
		return
	}

	err := c.pingWriter.Kill()
	c.pingWriter = nil
	if err != nil {
		err = fmt.Errorf("stop ping: error killing ping command: %v", err)
		c.LogError(err)
	}
}

func (c *websocketClient) processRequest(req *request) {
	if req.Option == START_PING {
		c.startPing(req)
	} else if req.Option == STOP_PING {
		c.stopPing()
	}
}
