package gweb

type ISessions interface {
	HasSession(sid string) bool
}
