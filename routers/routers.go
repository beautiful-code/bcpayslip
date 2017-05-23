package routers

import (
	"net/http"

	"bcpayslip/controllers"
	"bcpayslip/middlewares"
	"bcpayslip/urls"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

// GetRouter registers all routes for the application ...
func GetRouter() *pat.Router {
	// url paths imported from urls package
	common := pat.New()
	// static route
	common.PathPrefix(urls.StaticPath).Handler(
		http.StripPrefix(urls.StaticPath, http.FileServer(http.Dir("static"))))
	// media route
	common.PathPrefix(urls.MediaPath).Handler(
		http.StripPrefix(urls.MediaPath, http.FileServer(http.Dir("media"))))
	// common routes
	common.Get(urls.AuthcallbackPath, controllers.AuthCallbackController)
	common.Get(urls.AuthPath, controllers.AuthController)
	common.Get(urls.LogoutPath, controllers.LogoutController)
	// payslip routes
	payslip := pat.New()
	payslip.Get(urls.HomePath, controllers.PayslipController)
	payslip.Get(urls.PayslipPath, controllers.PayslipController)
	payslip.Post(urls.PayslipPath, controllers.PayslipController)
	payslip.NotFoundHandler = http.HandlerFunc(controllers.NotFoundController)
	common.PathPrefix(urls.HomePath).Handler(
		negroni.New(
			negroni.HandlerFunc(
				middlewares.GothLoginMiddleware),
			negroni.HandlerFunc(
				middlewares.SetUserMiddleware),
			negroni.Wrap(payslip),
		),
	)
	common.Get(urls.NotfoundPath, controllers.NotFoundController)
	common.NotFoundHandler = http.HandlerFunc(controllers.NotFoundController)
	common.Get(urls.RootPath, controllers.LoginController)
	return common
}
