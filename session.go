package gweb

import (
	. "github.com/eynstudio/gobreak"
)

type ISession interface {
}

type Sessions map[GUID]ISession
