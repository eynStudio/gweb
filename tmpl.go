package gweb

import (
	"html/template"
	"io"
)

type Tmpl struct {
	Templates *template.Template
}

func (p *Tmpl) Load() {
	p.Templates = template.Must(template.New("").Delims("[[", "]]").ParseGlob("views/*.*"))
}

func (p *Tmpl) Execute(wr io.Writer, name string, data interface{}) error {
	return p.Templates.ExecuteTemplate(wr, name+".tpl", data)
}
