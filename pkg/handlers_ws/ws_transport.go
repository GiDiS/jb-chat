package handlers_ws

import (
	"context"
	"github.com/GiDiS/jb-chat/pkg/events"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"sync"
	"time"
)

type WsHandler func(context.Context, *http.Request, *websocket.Conn)

type WsApi interface {
	RegisterHandlers(r *mux.Router)
	Serve(ctx context.Context)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Subscribe(bus events.Connector)
	Unsubscribe(bus events.Connector)
	HasConnection(conn string) bool
}

type wsTransport struct {
	eventsBus        events.Connector
	eventsResolver   events.Resolver
	upgrader         websocket.Upgrader
	logger           logger.Logger
	clients          map[*WsClient]string
	connections      map[string]*WsClient
	exit             chan struct{}
	clientRegister   chan *WsClient
	clientUnregister chan *WsClient
	clientIncome     chan WsIncomeMsg
	busRecv          chan events.Event
	busSend          chan events.Event
	mx               sync.Mutex
}

func NewWsTransport(eventsResolver events.Resolver, log logger.Logger) *wsTransport {
	ws := wsTransport{
		eventsResolver:   eventsResolver,
		logger:           log.WithField("transport", "ws"),
		exit:             make(chan struct{}),
		clients:          make(map[*WsClient]string),
		connections:      make(map[string]*WsClient),
		clientIncome:     make(chan WsIncomeMsg, 4096),
		clientRegister:   make(chan *WsClient),
		clientUnregister: make(chan *WsClient),
		busRecv:          make(chan events.Event, 4096),
		busSend:          make(chan events.Event, 4096),

		upgrader: websocket.Upgrader{
			ReadBufferSize:    81920,
			WriteBufferSize:   81920,
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
				//return false
			},
		},
	}
	return &ws
}

func (a *wsTransport) Subscribe(bus events.Connector) {
	a.eventsBus = bus
	a.eventsBus.Subscribe(a.onBusRecv)
}

func (a *wsTransport) Unsubscribe(_ events.Connector) {
	if a.eventsBus != nil {
		a.eventsBus.Unsubscribe(a.onBusRecv)
		a.eventsBus = nil
	}
}

func (a *wsTransport) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/ws", a.upgradeAndRun())
}

func (a *wsTransport) Serve(ctx context.Context) {
	//a.eventsBus.Subscribe(a.onBusRecv)
	a.exit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				a.onExit()
			case client := <-a.clientRegister:
				a.onConnect(client)
			case client := <-a.clientUnregister:
				a.onDisconnect(client)
			case event := <-a.busSend:
				a.onEmit(event)
			case event := <-a.busRecv:
				a.onBusRecv(event)
			case msg, _ := <-a.clientIncome:
				a.onIncomeMsg(msg)
			}
		}
	}()
}

func (a *wsTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.upgradeAndRun()(w, r)
}

func (a *wsTransport) HasConnection(conn string) bool {
	a.mx.Lock()
	defer a.mx.Unlock()
	if _, ok := a.connections[conn]; ok {
		return true
	}
	return false
}

func (a *wsTransport) upgradeAndRun() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := a.upgrader.Upgrade(w, r, nil)
		if err != nil {
			a.logger.WithError(err).Errorf("Upgrade failed: %+v", err)
			return
		}
		client := NewWsClient(a, conn, a.logger)
		client.Serve(r.Context())
		hostname, err := os.Hostname()
		if err != nil {
		}
		a.clientRegister <- client
		a.logger.Debugf("Registered: %s, %s", hostname, client.GetId())
	}
}

func (a *wsTransport) onExit() {
	a.Unsubscribe(a.eventsBus)
	if a.exit != nil {
		close(a.exit)
		a.exit = nil
	}
}

func (a *wsTransport) onConnect(client *WsClient) {
	a.clients[client] = client.GetId()
	a.connections[client.GetId()] = client
	a.logger.Debugf("Connected: %s", client.GetId())
	if ev, err := NewClientEvent(WsConnected, client); err != nil {
		a.logger.Errorf("New connected event failed: %v", err)
	} else if ev != nil {
		a.onEmit(*ev)
	}
}

func (a *wsTransport) onDisconnect(client *WsClient) {
	if _, ok := a.clients[client]; ok {
		delete(a.clients, client)
		delete(a.connections, client.GetId())

		client.MarkDisconnected()
		a.logger.Debugf("Disconnected: %s", client.GetId())
		if ev, err := NewClientEvent(WsDisconnected, client); err != nil {
			a.logger.Errorf("New connected event failed: %v", err)
		} else if ev != nil {
			a.onEmit(*ev)
		}
	}
}

func (a *wsTransport) onIncomeMsg(msg WsIncomeMsg) {
	a.logger.Debugf("Recv: %+v", string(msg.msg))
	if recvEvent, err := a.unmarshalEvent(msg); err != nil {
		// send error msg back
		_ = msg.client.conn.WriteJSON(recvEvent)
	} else {
		//a.logger.Debugf("Recv event: %+v", recvEvent)
		a.onEmit(recvEvent)
	}
}

func (a *wsTransport) busRecvHandler(event events.Event) {
	a.busRecv <- event
}

func (a *wsTransport) onBusRecv(event events.Event) {
	if event.To == nil {
		return
	}

	switch event.To.Type {
	case events.DestinationBroadcast:
		for client := range a.clients {
			a.onDirect(client, event)
		}
	case events.DestinationConnection:
		if client, ok := a.connections[event.To.Addr]; ok {
			a.onDirect(client, event)
		}
	case events.DestinationNoop:
		// do nothing
	default:
		a.logger.Debugf("Unhandled destination: [%w,%s]", event.To.Type, event.To.Addr)

	}

}

func (a *wsTransport) onDirect(client *WsClient, event events.Event) {
	if client == nil {
		return
	} else if client.send == nil {
		return
	}

	msg, err := a.marshalEvent(client, event)
	if err != nil {
		return
	}

	a.logger.Debugf("Send: %+v", string(msg.msg))
	client.send <- msg.msg
}

func (a *wsTransport) onEmit(event events.Event) {
	if a.eventsBus != nil {
		a.eventsBus.Emit(event)
	}
}

func (a *wsTransport) marshalEvent(_ *WsClient, event events.Event) (WsOutcomeMsg, error) {
	blob, err := a.eventsResolver.Marshal(event)
	if err != nil {
		a.logger.Errorf("Outcome msg marshal failed: %v", err)
		return WsOutcomeMsg{}, err
	}

	msg := WsOutcomeMsg{msg: blob, at: event.At}
	if msg.at.IsZero() {
		msg.at = time.Now()
	}
	return msg, nil
}

func (a *wsTransport) unmarshalEvent(msg WsIncomeMsg) (events.Event, error) {
	recvEvent, err := a.eventsResolver.Unmarshal(msg.msg)
	if err != nil {
		return recvEvent, err
	}

	if recvEvent.Id == "" {
		recvEvent.Id = uuid.NewV4().String()
	}

	err = events.WithOptions(
		&recvEvent,
		events.WithCtx(msg.ctx),
		events.WithFrom(events.NewDestConnection(msg.client.GetId())),
		events.WithAt(msg.at),
	)

	return recvEvent, err
}
