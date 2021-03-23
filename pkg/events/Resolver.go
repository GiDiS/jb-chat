package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Resolver interface {
	Resolve(string, []byte) (Event, error)
	Register(eventType Type, proto interface{}) error
	Proto(eventType Type) (interface{}, error)
	UnmarshalPayload(eventType Type, rawPayload []byte) (interface{}, error)
	Unmarshal(raw []byte) (Event, error)
	Marshal(Event) ([]byte, error)
}

type resolver struct {
	eventProtos map[Type]interface{}
}

var DefaultResolver = &resolver{eventProtos: make(map[Type]interface{})}

func (r *resolver) Resolve(rawEventType string, rawPayload []byte) (Event, error) {

	event := Event{Type: Type(rawEventType)}
	if _, ok := r.eventProtos[event.Type]; !ok {
		return event, ErrEventNotRegistered
	}
	payload, err := r.UnmarshalPayload(event.Type, rawPayload)
	if err != nil {
		return event, err
	}

	err = WithPayload(payload)(&event)
	return event, err
}

func (r *resolver) Proto(eventType Type) (interface{}, error) {
	if proto, ok := r.eventProtos[eventType]; !ok {
		return nil, ErrEventNotRegistered
	} else {
		return proto, nil
	}

}
func (r *resolver) Marshal(event Event) ([]byte, error) {
	blob, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("event marshal failed: %w", err)
	}
	return blob, nil
}

func (r *resolver) Unmarshal(raw []byte) (Event, error) {
	var rawMsg map[string]json.RawMessage

	if err := json.Unmarshal(raw, &rawMsg); err != nil {
		return Event{Type: InvalidFormatted, Payload: raw}, fmt.Errorf("event json unmarshal failed: %w", err)
	}

	rawEventType, ok := rawMsg["type"]
	if !ok {
		return Event{Type: InvalidType, Payload: raw}, ErrUnmarshalTypeUndefined
	}
	event := Event{Type: Type(strings.Trim(string(rawEventType), "\""))}

	if rawPayload, ok := rawMsg["payload"]; ok {
		payload, err := DefaultResolver.UnmarshalPayload(event.Type, rawPayload)
		if err != nil {
			return Event{Type: InvalidPayload, Payload: string(raw)}, err
		}
		event.Payload = payload
	}

	rawId, ok := rawMsg["id"]
	if ok && string(rawId) != "\"\"" {
		event.Id = strings.Trim(string(rawId), "\"")
	}

	return event, nil
}

func (r *resolver) UnmarshalPayload(eventType Type, rawPayload []byte) (interface{}, error) {

	payload, err := r.newEventPayload(eventType)
	if err != nil {
		return nil, fmt.Errorf("event json unmarshal create payload failed: %w", err)
	}
	if payload == nil {
		if rawPayload == nil || string(rawPayload) == "null" || string(rawPayload) == "\"\"" {
			return nil, nil
		} else {
			return nil, fmt.Errorf("event json unmarshal failed: payload must be empty")
		}
	}
	if rawPayload == nil {
		return nil, fmt.Errorf("event json unmarshal failed: payload must not be empty")
	}
	if err = json.Unmarshal(rawPayload, payload); err != nil {
		return nil, fmt.Errorf("event json unmarshal payload failed: %w", err)
	}
	payloadVal := reflect.ValueOf(payload)
	payloadVal = reflect.Indirect(payloadVal)
	return payloadVal.Interface(), nil
}

func (r *resolver) Register(eventType Type, proto interface{}) error {
	if _, ok := r.eventProtos[eventType]; ok {
		return ErrEventAlreadyRegistered
	}

	protoVal := reflect.ValueOf(proto)
	if protoVal.Kind() == reflect.Ptr {
		return ErrProtoMustByValue
	}

	r.eventProtos[eventType] = proto
	return nil
}

func (r *resolver) newEventPayload(t Type) (interface{}, error) {
	proto, ok := r.eventProtos[t]
	if !ok {
		return nil, ErrEventNotRegistered
	}

	if proto == nil {
		return nil, nil
	}

	v := reflect.ValueOf(proto)
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Struct:
		eventPayload := reflect.New(v.Type())
		return eventPayload.Interface(), nil
	case reflect.String:
		return "", nil
	default:
		return nil, errors.New("invalid type")
	}
}
