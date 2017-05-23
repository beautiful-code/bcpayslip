package controllers

import (
	"net/http"
	"strings"
	"text/template"

	"bcpayslip/store"
	"bcpayslip/templates"
	"bcpayslip/urls"
	"bcpayslip/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// LoginController login page controller ...
func LoginController(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.GetValidSession(req)
	if session.Values["userid"] != nil {
		context.Set(req, "userid", session.Values["userid"])
		http.Redirect(res, req, urls.HomePath, http.StatusSeeOther)
	} else {
		t, _ := template.ParseFiles(templates.LoginTemplate)
		t.Execute(res, nil)
	}
}

// LogoutController delete the cookie and redirect ...
func LogoutController(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.GetValidSession(req)
	session.Options = &sessions.Options{Path: urls.RootPath, MaxAge: -1}
	session.Save(req, res)
	http.Redirect(res, req, urls.RootPath, http.StatusTemporaryRedirect)
}

// AuthController start authentication process using gothic ...
func AuthController(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

// AuthCallbackController goth callback controller to complete user auth and create user ...
func AuthCallbackController(res http.ResponseWriter, req *http.Request) {
	var gothUser goth.User
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		gothic.BeginAuthHandler(res, req)
	}
	emailDomain := strings.Join(strings.Split(gothUser.Email, "@")[1:], "a")
	if emailDomain != "beautifulcode.in" {
		session, _ := utils.GetValidSession(req)
		session.Options = &sessions.Options{Path: urls.RootPath, MaxAge: -1}
		session.Save(req, res)
		queryParam := "?m=Invalid account, use BC account"
		http.Redirect(res, req, urls.RootPath+queryParam, http.StatusSeeOther)
	}
	session, _ := utils.GetValidSession(req)
	session.Values["userid"] = gothUser.UserID
	session.Save(req, res)
	store.SaveUser(
		gothUser.UserID, gothUser.FirstName, gothUser.LastName,
		gothUser.Email, gothUser.AccessToken, gothUser.AvatarURL,
	)
	http.Redirect(res, req, urls.HomePath, http.StatusSeeOther)
}
