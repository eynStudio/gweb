package gweb

import (
	. "github.com/eynstudio/gobreak"
)

type ISession interface {
	ID() string
	Uid() GUID
}

type ISessions interface {
	HasSession(sid string) (bool, error)
	DelSession(sid string) error
	GetSession(sid string) (ISession, error)
}
