package translation

import (
	"bytes"
	"encoding"
	"math/rand"
	"strings"
	gotemplate "text/template"
)

type template struct {
	tmpl *gotemplate.Template
	src  string
}

func newTemplate(src string) (*template, error) {
	if src == "" {
		return new(template), nil
	}

	var tmpl template
	err := tmpl.parseTemplate(src)
	return &tmpl, err
}

func mustNewTemplate(src string) *template {
	t, err := newTemplate(src)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *template) String() string {
	return t.src
}

func (t *template) Execute(args interface{}) string {
	if t.tmpl == nil {
		return t.src
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, args); err != nil {
		return err.Error()
	}
	return buf.String()
}

func (t *template) MarshalText() ([]byte, error) {
	return []byte(t.src), nil
}

func (t *template) UnmarshalText(src []byte) error {
	return t.parseTemplate(string(src))
}

func (t *template) parseTemplate(src string) (err error) {
	t.src = src
	if strings.Contains(src, "{{") {
		t.tmpl, err = gotemplate.New(src).Funcs(t.addTemplateFunctions()).Parse(src)
	}
	return
}

// Adds needed custom template functions
func (t *template) addTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Random": func(values ...string) string {
			return values[rand.Intn(len(values))]
		},
		"Username": func(user map[string]string) string {
			if _, ok := user["name"]; !ok {
				return "<no value>"
			}

			return user["name"]
		},
		"Subtract": func(a, b int) int {
			return a - b
		},
	}
}

var _ = encoding.TextMarshaler(&template{})
var _ = encoding.TextUnmarshaler(&template{})
