package gopher

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jimmykuu/wtforms"
)

const (
	BASE  = "base.html"
	ADMIN = "admin/base.html"
)

var funcMaps = template.FuncMap{
	"input": func(form wtforms.Form, fieldStr string) template.HTML {
		field, err := form.Field(fieldStr)
		if err != nil {
			panic(err)
		}

		errorClass := ""
		errorMessage := ""
		if field.HasErrors() {
			errorClass = "error "
			errorMessage = `<div class="ui red pointing above ui label">` + strings.Join(field.Errors(), ", ") + `</div>`
		}
		format := `<div class="%sfield">
			    %s
			    %s
				%s
		    </div>`

		return template.HTML(
			fmt.Sprintf(format,
				errorClass,
				field.RenderLabel(),
				field.RenderInput(),
				errorMessage,
			))
	},
	"loadtimes": func(startTime time.Time) string {
		return fmt.Sprintf("%dms", time.Now().Sub(startTime)/1000000)
	},
}

// 解析模板
func parseTemplate(file, baseFile string, data map[string]interface{}) []byte {
	var buf bytes.Buffer
	t := template.New(file).Funcs(funcMaps)
	baseBytes, err := ioutil.ReadFile("templates/" + baseFile)
	if err != nil {
		panic(err)
	}
	t, err = t.Parse(string(baseBytes))
	if err != nil {
		panic(err)
	}
	t, err = t.ParseFiles("templates/" + file)
	if err != nil {
		panic(err)
	}
	err = t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}