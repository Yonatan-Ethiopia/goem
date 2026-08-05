package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/justinas/alice"
)

func (app *application) routes() *httprouter.Router{
    
    baseMiddleware := alice.New(app.recoverPanic, secureHeaders, app.sessionManager.LoadAndSave, noSurf,app.authenticate)
    protected := baseMiddleware.Append(app.requireAuthentication)
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    
    router := httprouter.New()
    
    router.HandlerFunc(http.MethodGet, "/ping", ping)
    
    router.Handler(http.MethodGet, "/", baseMiddleware.ThenFunc(http.HandlerFunc(app.homePage)))
    router.Handler(http.MethodGet, "/still/:id", baseMiddleware.ThenFunc(http.HandlerFunc(app.boxViewGet)))
    router.Handler(http.MethodGet, "/create", protected.ThenFunc(http.HandlerFunc(app.formCreateGet)))
    router.Handler(http.MethodPost, "/create", protected.ThenFunc(http.HandlerFunc(app.formCreatePost)))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

         
    router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignUpGet))))
    router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignUpPost))))
    
    router.Handler(http.MethodGet, "/user/login", baseMiddleware.ThenFunc(http.HandlerFunc(app.userLogInGet)))
    router.Handler(http.MethodPost, "/user/login", baseMiddleware.ThenFunc(http.HandlerFunc(app.userLogInPost)))
    router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(http.HandlerFunc(app.userLogOutPost)))
    
    return router
}
