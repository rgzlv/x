package util

import (
	htmlT "html/template"
	"net/http"
	"path/filepath"
	textT "text/template"
)

type TmplData struct {
	URLPath string
	Auth    Auth
	Data    any
	ErrMsg  string
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, filename string, tmplDir string, data any) error {
	var err error

	tmpl, err := htmlT.ParseGlob(tmplDir + "/*.tmpl.html")
	if err != nil {
		LogError(err.Error())
		return err
	}

	_, err = tmpl.ParseFiles(filename)
	if err != nil {
		LogError(err.Error())
		return err
	}

	filename = filepath.Base(filename)
	err = tmpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		LogError(err.Error())
		return err
	}

	return nil
}

func ExecuteTemplateHTML(w http.ResponseWriter, r *http.Request, filename string, tmplDir string, data any) error {
	var err error

	tmpl, err := textT.ParseGlob(tmplDir + "/*.tmpl.html")
	if err != nil {
		return err
	}

	_, err = tmpl.ParseFiles(filename)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		return err
	}

	return nil
}

func ExecuteTemplateError(w http.ResponseWriter, r *http.Request, tmplDir string, tmplData *TmplData, msg string) {
	var err error

	tmplData.ErrMsg = msg
	LogError(tmplData.ErrMsg)
	err = ExecuteTemplate(w, r, "error.html", tmplDir, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ExecuteTemplateLoginWithError(w http.ResponseWriter, r *http.Request, tmplData *TmplData, tmplDir string, msg error) {
	var err error

	tmplData.Auth.Status = ASError
	tmplData.Auth.Error = msg.Error()

	LogError(msg.Error())
	err = ExecuteTemplate(w, r, "login.html", tmplDir, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
