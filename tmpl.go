package gweb

import (
	"html/template"
	"io"
)

type Tmpl struct {
	Templates *template.Template
}

func html(v string) template.HTML { return template.HTML(v) }

func (p *Tmpl) Load() {
	p.Templates = template.Must(template.New("").Funcs(template.FuncMap{
		"html": html,
	}).Delims("[[", "]]").ParseGlob("views/*.*"))
}

func (p *Tmpl) Execute(wr io.Writer, name string, data interface{}) error {
	return p.Templates.ExecuteTemplate(wr, name+".tpl", data)
}
