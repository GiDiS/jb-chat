package auth

import "github.com/GiDiS/jb-chat/pkg/events"

const (
	AuthMe         events.Type = "auth.me"
	AuthMeInfo     events.Type = "auth.me-info"
	AuthRegister   events.Type = "auth.register"
	AuthRegistered events.Type = "auth.registered"
	AuthRequired   events.Type = "auth.required"
	AuthSignIn     events.Type = "auth.sign-in"
	AuthSignedIn   events.Type = "auth.signed-in"
	AuthSignOut    events.Type = "auth.sign-out"
	AuthSignedOut  events.Type = "auth.signed-out"
)
