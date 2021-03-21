package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
	"time"
)

type Type string

const InvalidType Type = "sys.invalid_payload"
const InvalidFormatted Type = "sys.invalid_type"
const InvalidPayload Type = "sys.invalid_payload"

type DestinationType int

const (
	DestinationNoop DestinationType = iota
	DestinationChannel
	DestinationUser
	DestinationSession
	DestinationConnection
	DestinationBroadcast
)

func (t Type) String() string {
	return string(t)
}

type Event struct {
	Id      string          `json:"id"`
	Type    Type            `json:"type"`
	Payload interface{}     `json:"payload,omitempty"`
	From    *Destination    `json:"from,omitempty"`
	To      *Destination    `json:"to,omitempty"`
	At      time.Time       `json:"at,omitempty"`
	Prev    string          `json:"prev,omitempty"`
	Ctx     context.Context `json:"-"`
	proto   interface{}
}

type EventOptionSetter func(e *Event) error

type Destination struct {
	Type DestinationType `json:"type"`
	Addr string          `json:"addr,omitempty"`
}

func NewEvent(eventType Type, opts ...EventOptionSetter) (Event, error) {
	event := Event{Type: eventType, At: time.Now()}
	proto, err := DefaultResolver.Proto(eventType)
	if err != nil {
		return event, err
	}
	event.proto = proto

	if err := WithOptions(&event, opts...); err != nil {
		return event, err
	}

	if event.Id == "" {
		event.GetId()
	}

	return event, nil
}

func NewDestNoop() *Destination {
	return &Destination{Type: DestinationNoop}
}

func NewDestBroadcast() *Destination {
	return &Destination{Type: DestinationBroadcast}
}

func NewDestUser(uid models.Uid) *Destination {
	return &Destination{Type: DestinationUser, Addr: uid.String()}
}

func NewDestSession(sid models.SessionId) *Destination {
	return &Destination{Type: DestinationSession, Addr: sid.String()}
}

func NewDestChannel(cid models.ChannelId) *Destination {
	return &Destination{Type: DestinationChannel, Addr: cid.String()}
}

func NewDestConnection(connName string) *Destination {
	return &Destination{Type: DestinationConnection, Addr: connName}
}

func (e *Event) GetId() string {
	if e.Id == "" {
		e.Id = uuid.NewV4().String()
	}
	return e.Id
}

// GetSid - gets associated session id
func (e *Event) GetSid() string {
	if e == nil || e.Ctx == nil {
		return ""
	}
	conId, ok := e.Ctx.Value("connection").(string)
	if !ok {
		return ""
	}
	return conId
}

// GetUid - gets associated user id
func (e *Event) GetUid() models.Uid {
	if e == nil || e.Ctx == nil {
		return models.NoUser
	}
	uid, ok := e.Ctx.Value("uid").(models.Uid)
	if !ok {
		return models.NoUser
	} else {
		return uid
	}
}

// SetUid - sets associated user id
func (e *Event) SetUid(uid models.Uid) {
	if e == nil {
		return
	}
	if e.Ctx == nil {
		e.Ctx = context.Background()
	}
	e.Ctx = context.WithValue(e.Ctx, "uid", uid)
}

func (e *Event) UnmarshalJSON(blob []byte) error {

	var raw map[string]json.RawMessage

	if err := json.Unmarshal(blob, &raw); err != nil {
		return fmt.Errorf("event json unmarshal failed: %w", err)
	}

	rawEventType, ok := raw["type"]
	if !ok {
		return ErrUnmarshalTypeUndefined
	}
	*e = Event{Type: Type(strings.Trim(string(rawEventType), "\""))}

	rawPayload, _ := raw["payload"]
	payload, err := DefaultResolver.UnmarshalPayload(e.Type, rawPayload)
	if err != nil {
		return err
	}
	e.Payload = payload

	return nil
}

func WithOptions(event *Event, opts ...EventOptionSetter) error {
	for _, opt := range opts {
		if err := opt(event); err != nil {
			return err
		}
	}
	return nil
}

func WithPayload(payload interface{}) EventOptionSetter {
	return func(e *Event) error {
		if e == nil {
			return errors.New("nil event")
		}
		payloadType := reflect.TypeOf(payload)
		protoType := reflect.TypeOf(e.proto)
		if payloadType != protoType {
			return errors.New("wrong payload type")
		}
		e.Payload = payload

		if e.proto != nil && e.Payload == nil {
			return errors.New("must be only one payload")
		}
		return nil
	}
}

func WithFrom(dest *Destination) EventOptionSetter {
	return func(e *Event) error {
		e.From = dest
		return nil
	}
}

func WithTo(dest *Destination) EventOptionSetter {
	return func(e *Event) error {
		e.To = dest
		return nil
	}
}

func WithPrev(event *Event) EventOptionSetter {
	return func(e *Event) error {
		if event != nil {
			e.Prev = event.Id
		} else {
			e.Prev = ""
		}
		return nil
	}
}

func WithAt(at time.Time) EventOptionSetter {
	return func(e *Event) error {
		e.At = at
		return nil
	}
}

func WithCtx(ctx context.Context) EventOptionSetter {
	return func(e *Event) error {
		e.Ctx = ctx
		return nil
	}
}

type ResultStatus struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func (r *ResultStatus) IsFailed() bool {
	return !r.Ok
}

func (r *ResultStatus) SetSuccess(msg string) {
	r.Ok = true
	r.Message = msg
}

func (r *ResultStatus) SetFailed(err error) {
	r.Ok = false
	r.Message = err.Error()
}
