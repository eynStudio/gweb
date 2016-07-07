package gweb

type ISessions interface {
	HasSess(sid string) bool
	DelSess(sid string) error
	GetSessUid(sid string) (string, error)
}
