package main

import (
        "net/http"
        "context"
        "fmt"
        "github.com/justinas/nosurf"
        "log"
)

func noSurf( next http.Handler) http.Handler{
    csrfHandler := nosurf.New(next)
    
    csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("NOSURF FAILED: reason=%v", nosurf.Reason(r))
        http.Error(w, "Bad Request", http.StatusBadRequest)
    }))
    csrfHandler.SetBaseCookie(http.Cookie{
        HttpOnly: true,
        Path: "/",
        Secure: true,
    })
    return csrfHandler
}

func secureHeaders( next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "deny")
        w.Header().Set("X-XSS-Protection", "0")
        
        next.ServeHTTP(w,r)
    })
}

func (app *application) recoverPanic(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        
        defer func(){
            if err := recover(); err != nil{
                w.Header().Set("Connection", "close")
                
                app.serverError(w, fmt.Errorf("%s", err))
            }
        }()
        
        next.ServeHTTP(w, r)
        
    })
}

func (app *application) authenticate(next http.Handler) http.Handler{
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
    if id == 0{
        next.ServeHTTP(w, r)
        return
    }
    exists, err := app.userconn.Exists(id)
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    if exists {
        ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
        r = r.WithContext(ctx)
    }
    
    next.ServeHTTP(w, r)
    })
}

func (app *application) requireAuthentication( next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
            if !app.isAuthenticated(r){
                http.Redirect(w,r,"/user/login", http.StatusSeeOther)
                return
            }
            
            w.Header().Add("Cache-Control", "no-store")
            
            next.ServeHTTP(w,r)
        })
}
