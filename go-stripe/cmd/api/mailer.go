package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var emailTempateFS embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data any) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTempateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)

	t, err = template.New("email-plain").ParseFS(emailTempateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}
