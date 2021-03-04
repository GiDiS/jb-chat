package handlers_ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"jb_chat/pkg/models"
	"log"
	"sync/atomic"
	"time"
)

const (
	readWait       = 2 * time.Second
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 51200
)

var wsClientSeq int32 = 0

type WsClient struct {
	id        int
	transport *wsTransport
	ctx       context.Context
	conn      *websocket.Conn
	uid       models.Uid
	send      chan []byte
}

type WsIncomeMsg struct {
	client *WsClient
	ctx    context.Context
	msg    []byte
	at     time.Time
}

type WsOutcomeMsg struct {
	ctx context.Context
	msg []byte
	at  time.Time
}

func NewWsClient(transport *wsTransport, conn *websocket.Conn) *WsClient {

	id := atomic.AddInt32(&wsClientSeq, 1)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "transport", "websocket")
	ctx = context.WithValue(ctx, "connection", id)
	ctx = context.WithValue(ctx, "remote_addr", conn.RemoteAddr())
	ctx = context.WithValue(ctx, "local_addr", conn.LocalAddr())

	return &WsClient{
		id:        int(id),
		transport: transport,
		conn:      conn,
		ctx:       ctx,
		send:      make(chan []byte, 256),
	}
}

func (c *WsClient) Serve(ctx context.Context) {
	go c.readPump(ctx)
	go c.writePump(ctx)
}

func (c *WsClient) GetId() string {
	if c.id == 0 {
		id := atomic.AddInt32(&wsClientSeq, 1)
		c.id = int(id)
	}
	return fmt.Sprintf("ws-%d", c.id)
}

func (c *WsClient) readPump(ctx context.Context) {
	defer func() {
		c.transport.clientUnregister <- c
		_ = c.conn.Close()
	}()

	//c.conn.SetReadLimit(maxMessageSize)
	//_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error {
	//	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//	return nil
	//})

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Connection error: %v", err)
			} else {
				log.Printf("Unknown error: %v", err)
			}
			break
		}
		switch messageType {
		case websocket.TextMessage, websocket.BinaryMessage:
			log.Printf("Rcv: %s", string(message))
			c.transport.clientIncome <- c.newIncome(message)
		default:
			continue
		}
		//select {
		//case <-ctx.Done():
		//	break
		//}
	}
}

func (c *WsClient) newIncome(msg []byte) WsIncomeMsg {
	return WsIncomeMsg{
		client: c,
		ctx:    c.ctx,
		msg:    msg,
		at:     time.Now(),
	}
}

func (c *WsClient) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		//case <-ctx.Done():
		//	return
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Offload all messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte("\n"))
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
