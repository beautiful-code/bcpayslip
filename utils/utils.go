package utils

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"bcpayslip/models"
	"bcpayslip/store"
	"bcpayslip/templates"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
)

// GetValidSession Returns a valid authenticated user session ...
func GetValidSession(req *http.Request) (*sessions.Session, error) {
	sessStore := sessions.NewCookieStore([]byte(os.Getenv("bc_app_key")))
	return sessions.Store.Get(sessStore, req, "gplus_gothic_session")
}

// CustomTemplateExecute Append common templates and data structs and execute template ...
func CustomTemplateExecute(res http.ResponseWriter, req *http.Request, templateName string, data map[string]interface{}) {
	t, _ := template.ParseFiles(templates.BaseTemplate, templateName)
	if len(data) == 0 {
		data = make(map[string]interface{})
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
	} else {
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
	}
	if err := t.Execute(res, data); err != nil {
		log.Println(err)
	}
}

// AddParamsToURL Add params to url using a splice of models.kwargs struct ...
func AddParamsToURL(url string, args []models.Kwargs) string {
	for _, arg := range args {
		url = strings.Replace(url, "{"+arg.Key+"}", arg.Value, 1)
	}
	return url
}
