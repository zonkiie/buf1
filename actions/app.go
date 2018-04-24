package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"buf1/models"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

var logger buffalo.Logger = nil

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	logger = buffalo.NewLogger("debug")

	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_buf1_session",
			LooseSlash:  true,
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.ANY("/models/list", ModelListHandler)
		//logger.Infof("App: %v", app)
		//logger.Infof("Session: %v", app.Context.Session().Get("current_user_id"))
		app.ANY("/user/add", UserAdd)
		app.ANY("/user/addcommit", UserAddcommit)
		app.ANY("/user/get", UserGet)
		app.ANY("/user/edit", UserEdit)
		app.POST("/user/remove", UserRemove)
		app.POST("/user/removeconfirm", UserRemoveconfirm)
		app.ANY("/user/editcommit", UserEditcommit)
		app.ANY("/user/login", UserLogin)
		app.ANY("/user/logincommit", UserLogincommit)
		app.ANY("/user/logout", UserLogout)
		//app.ANY("/book/add", BookAdd)

		app.ANY("/book/add", BookAdd)
		app.ANY("/book/edit", BookEdit)
		app.ANY("/book/editcommit", BookEditcommit)
		app.ANY("/book/addcommit", BookAddcommit)
		app.POST("/book/removecommit", BookRemovecommit)
		app.POST("/book/remove", BookRemove)
		app.ANY("/book/list", BookList)
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
