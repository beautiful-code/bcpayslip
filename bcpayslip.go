package main

import (
	// system local third-party

	"net/http"
	"os"

	"bcpayslip/routers"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"github.com/urfave/negroni"
)

func init() {
	// goth package cookie store initialization
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("bc_app_key")))
}

func main() {
	StartMyApp()
}

func StartMyApp() {
	// use goth provider for authentication
	if os.Getenv("bc_env") == "development" {
		goth.UseProviders(
			gplus.New(
				os.Getenv("bc_intranet_client_id"),
				os.Getenv("bc_intranet_client_secret"),
				os.Getenv("bc_host")+":"+os.Getenv("PORT")+"/auth/gplus/callback",
			),
		)
	} else {
		goth.UseProviders(
			gplus.New(
				os.Getenv("bc_intranet_client_id"),
				os.Getenv("bc_intranet_client_secret"),
				os.Getenv("bc_host")+"/auth/gplus/callback",
			),
		)
	}
	// get pat router from routers package
	p := routers.GetRouter()
	// use negroni handler
	n := negroni.Classic()
	n.UseHandler(p)
	// run on 3001 and using gin(repl) on 3000
	var port string
	if os.Getenv("bc_env") == "development" {
		port = "3001"
	} else {
		port = os.Getenv("PORT")
	}
	http.ListenAndServe(":"+port, n)
}
