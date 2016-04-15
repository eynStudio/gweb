package gweb

type Cfg struct {
	Port       int
	Tls        bool
	CertFile   string
	KeyFile    string
	ServeFiles []string
	useTmpl    bool
}
