package gweb

type Cfg struct {
	Http       bool
	HttpPort   int
	Https      bool
	HttpsPort  int
	CertFile   string
	KeyFile    string
	ServeFiles []string
	UseTmpl    bool
}
