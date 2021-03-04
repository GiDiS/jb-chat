package events

import (
	"context"
	"jb_chat/pkg/logger"
)

type Dispatcher interface {
	AddTransport(tr Transport, local bool)
	Emit(Event) error
	Notify(e Event) error
	On(t Type, h EventHandler)
}

type dispatcher struct {
	transports map[Connector]bool
	listeners  map[Type][]EventHandler
	logger     logger.Logger
}

type EventHandler func(e Event) error

func NewDispatcher(
	logger logger.Logger,
) *dispatcher {
	d := dispatcher{
		transports: make(map[Connector]bool),
		listeners:  make(map[Type][]EventHandler),
		logger:     logger,
	}
	return &d
}

func (d *dispatcher) AddTransport(tr Transport, local bool) {
	busSize := 4096
	bus := NewConnector(busSize, busSize)
	bus.Serve(context.Background())
	tr.Subscribe(bus)
	d.transports[bus] = local
	go func() {
		if bus.sendChan != nil {
			for e := range bus.sendChan {
				d.Emit(e)
			}
		}
	}()
}

func (d *dispatcher) Emit(e Event) error {
	ls, ok := d.listeners[e.Type]
	if !ok {
		return nil
	}
	for _, l := range ls {
		err := l(e)
		if err != nil {
			d.logger.Error(err)
			return err
		}

	}
	return nil
}

func (d *dispatcher) Notify(e Event) error {
	for bus := range d.transports {
		bus.Notify(e)
	}
	return nil
}

func (d *dispatcher) On(t Type, h EventHandler) {
	if _, ok := d.listeners[t]; !ok {
		d.listeners[t] = make([]EventHandler, 0)
	}
	d.listeners[t] = append(d.listeners[t], h)
}
