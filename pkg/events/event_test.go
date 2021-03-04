package events

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const EmptyEvent Type = "test.empty"
const StructEvent Type = "test.struct"

type StructEventPayload struct{}

func TestNewEvent(t *testing.T) {
	initResolver(t)

	event, err := NewEvent(EmptyEvent)
	assert.NoError(t, err)
	assert.Equal(t, EmptyEvent, event.Type)
	assert.Nil(t, event.Payload)

	event, err = NewEvent(EmptyEvent, WithPayload("must fail"))
	assert.Error(t, err)
	assert.Equal(t, EmptyEvent, event.Type)
	assert.Nil(t, event.Payload)

	event, err = NewEvent(StructEvent, WithPayload(StructEventPayload{}))
	assert.NoError(t, err)
	assert.Equal(t, StructEvent, event.Type)
	assert.Equal(t, interface{}(StructEventPayload{}), event.Payload)

	event, err = NewEvent(StructEvent, WithPayload("must fail"))
	assert.EqualError(t, err, "wrong payload type")
	assert.Equal(t, StructEvent, event.Type)
	assert.Equal(t, nil, event.Payload)

}

func TestEvent_MarshalJSON(t *testing.T) {
	initResolver(t)
	event, err := NewEvent(EmptyEvent)
	assert.NoError(t, err)

	blob, err := json.Marshal(event)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"type":"test.empty"}`), blob)

	event = Event{Type: StructEvent, Payload: interface{}(StructEventPayload{})}

	blob, err = json.Marshal(event)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"test.struct","payload":{}}`, string(blob))
}

func TestEvent_UnmarshalJSON(t *testing.T) {
	initResolver(t)

	event := Event{}
	err := json.Unmarshal([]byte(`{"type":"test.empty"}`), &event)
	assert.NoError(t, err)

	event = Event{}
	err = json.Unmarshal([]byte(`{"type":"test.empty", "payload": null}`), &event)
	assert.NoError(t, err)

	event = Event{}
	err = json.Unmarshal([]byte(`{"type":"test.empty", "payload": ""}`), &event)
	assert.NoError(t, err)

	event = Event{}
	err = json.Unmarshal([]byte(`{"type":"test.empty", "payload": "foo"}`), &event)
	assert.EqualError(t, err, "event json unmarshal failed: payload must be empty")

	event = Event{}
	err = json.Unmarshal([]byte(`{"type":"test.struct","payload": {}}`), &event)
	assert.NoError(t, err)
	assert.Equal(t, interface{}(StructEventPayload{}), event.Payload)

	event = Event{}
	err = json.Unmarshal([]byte(`{"type":"test.struct"}`), &event)
	assert.EqualError(t, err, "event json unmarshal failed: payload must not be empty")
}

func initResolver(t *testing.T) {
	_ = DefaultResolver.Register(EmptyEvent, nil)
	//assert.NoError(t, err)
	_ = DefaultResolver.Register(StructEvent, StructEventPayload{})
	//assert.NoError(t, err)
}
