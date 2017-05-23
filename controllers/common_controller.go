package controllers

import (
	"net/http"
	"text/template"

	"bcpayslip/templates"
)

// NotFoundController 404 ...
func NotFoundController(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles(templates.NotfoundTemplate)
	t.Execute(res, nil)
}
