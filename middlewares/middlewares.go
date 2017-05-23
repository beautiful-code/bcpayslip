package middlewares

import (
	"net/http"

	"bcpayslip/urls"
	"bcpayslip/utils"

	"github.com/gorilla/context"
)

// GothLoginMiddleware Retreiving session, redirecting if no session found ...
func GothLoginMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	session, _ := utils.GetValidSession(req)
	if session.Values["gplus"] == nil {
		http.Redirect(res, req, urls.RootPath, http.StatusSeeOther)
	}
	if session.Values["userid"] != nil {
		context.Set(req, "userid", session.Values["userid"])
	} else {
		http.Redirect(res, req, urls.LogoutPath, http.StatusSeeOther)
	}
	next(res, req)
}

// SetUserMiddleware Appending the user id to every request and redirecting accordinly if no profile found ...
func SetUserMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	session, _ := utils.GetValidSession(req)
	if session.Values["userid"] != nil {
		context.Set(req, "userid", session.Values["userid"])
	} else {
		http.Redirect(res, req, urls.LogoutPath, http.StatusSeeOther)
	}
	next(res, req)
}
