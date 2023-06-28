package handlers

import (
	"oserver/structs"
	"time"
)



// TODO : replace sessions map with database or redis

type Session struct {
	userId 			string
	expiry 			time.Time
}

var Sessions = map[string]Session{}

func (s Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Username 			string					`json:"username"`
	Email 				string					`json:"email"`
	Password 			string					`json:"password"`
}


func GenerateSessionExpiry(t time.Duration) time.Time {
	return time.Now().Add(t)
}

func GenerateSessionToken() string {
	return "session-" + structs.GenerateID()
}