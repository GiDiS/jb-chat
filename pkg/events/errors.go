package events

import "errors"

var ErrEventAlreadyRegistered = errors.New("event type already registered")
var ErrEventNotRegistered = errors.New("event type not registered")
var ErrProtoMustByValue = errors.New("event proto must not be ptr")
var ErrUnmarshalTypeUndefined = errors.New("event json unmarshal failed: type undefined")
