package template

import (
	"bytes"
	"io"
	"text/template"
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/lmacrc/weather/pkg/sanitize"
)

type Template struct {
	inner *template.Template
}

// New allocates a new, undefined template with the given name.
func New(name string) *Template {
	t := template.New(name)
	t.Funcs(template.FuncMap{
		"strftime": strftimeFn,
	})

	return &Template{t}
}

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	var buf bytes.Buffer
	err := t.inner.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return err
	}

	s := sanitize.BaseName(buf.String())
	_, err = wr.Write([]byte(s))
	return err
}

func (t *Template) Execute(wr io.Writer, data interface{}) error {
	var buf bytes.Buffer
	err := t.inner.Execute(&buf, data)
	if err != nil {
		return err
	}

	s := sanitize.BaseName(buf.String())
	_, err = wr.Write([]byte(s))
	return err
}

func (t *Template) Parse(text string) (*Template, error) {
	_, err := t.inner.Parse(text)
	return t, err
}

func strftimeFn(f string, t time.Time) (string, error) {
	o, err := strftime.New(f)
	if err != nil {
		return "", err
	}
	return o.FormatString(t), nil
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//	var t = template.Must(template.New("name").Parse("text"))
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
