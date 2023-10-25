package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, app.noSurf, app.authenticate)
	authenticatedDynamic := dynamicMiddleware.Append(app.requireAuthentication)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	// User routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", authenticatedDynamic.ThenFunc(app.logoutUser))

	// Snippet routes
	mux.Get("/snippet/create", authenticatedDynamic.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", authenticatedDynamic.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", authenticatedDynamic.ThenFunc(app.showSnippet))

	fileHandler := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", disalowFolderListing(app, fileHandler)))

	return standardMiddleware.Then(mux)
}
