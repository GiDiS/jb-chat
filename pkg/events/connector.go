package events

import (
	"context"
	"sync"
)

type Connector interface {
	Subscribe(receiver func(event Event))
	Unsubscribe(receiver func(event Event))
	Emit(event Event)
	Notify(event Event)
	//Recv() chan<- Event
}

type BusStats interface {
	RecvLen() int
	SendLen() int
}

type connector struct {
	recvChan chan Event
	sendChan chan Event
	onRecv   func(Event)
	mx       sync.Mutex
}

func NewConnector(recvSize, sendSize int) *connector {
	bus := connector{}
	if recvSize >= 0 {
		bus.recvChan = make(chan Event, recvSize)
	}

	if sendSize >= 0 {
		bus.sendChan = make(chan Event, sendSize)
	}

	return &bus
}

func NewSendBus(size int) *connector {
	return NewConnector(-1, size)
}

func NewRecvBus(size int) *connector {
	return NewConnector(size, -1)
}

func (c *connector) Serve(ctx context.Context) {
	go func() {
		for {
			select {
			case income, ok := <-c.recvChan:
				if !ok {
					return
				} else if c.onRecv != nil {
					c.onRecv(income)
				}
			case <-ctx.Done():
				c.mx.Lock()
				if c.sendChan != nil {
					close(c.sendChan)
				}
				c.sendChan = nil
				c.mx.Unlock()
				return
			}
		}
	}()
}
func (c *connector) Subscribe(receiver func(event Event)) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.onRecv = receiver
}

func (c *connector) Unsubscribe(receiver func(event Event)) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.onRecv = nil
}

func (c *connector) Emit(event Event) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.sendChan != nil {
		c.sendChan <- event
	}
}

func (c *connector) Notify(event Event) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.recvChan != nil {
		c.recvChan <- event
	}
}

func (c *connector) RecvLen() int {
	if c.recvChan != nil {
		return len(c.recvChan)
	}
	return -1
}

func (c *connector) SendLen() int {
	if c.sendChan != nil {
		return len(c.sendChan)
	}
	return -1
}
