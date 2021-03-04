package events

type Transport interface {
	HasConnection(conn string) bool
	Subscribe(bus Connector)
	Unsubscribe(bus Connector)
}
