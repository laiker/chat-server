package model

import "time"

type AnonymousUser struct {
	id        int64
	login     string
	createdAt time.Time
}

func NewAnonymousUser(login string) *AnonymousUser {
	temporaryID := -1 * time.Now().UnixNano()

	return &AnonymousUser{
		id:        temporaryID,
		login:     login,
		createdAt: time.Time{},
	}
}

func (a *AnonymousUser) GetID() int64 {
	return a.id
}

func (a *AnonymousUser) GetLogin() string {
	return a.login
}

func (a *AnonymousUser) SetLogin(login string) {
	a.login = login
}

func (a *AnonymousUser) GetCreatedAt() time.Time {
	return a.createdAt
}
