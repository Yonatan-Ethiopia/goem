package main

import (
        "net/http"
        "html/template"
        "log"
        "flag"
        "os"
        "time"
        "backgo/internal/models"
        
        _ "github.com/jackc/pgx/v5/stdlib"
        "github.com/go-playground/form/v4"
        "github.com/alexedwards/scs/v2"
        "github.com/alexedwards/scs/postgresstore"
        
    )
    
type application struct{
    errLog *log.Logger
    infoLog *log.Logger
    dbconn models.BoxInterface
    userconn models.UserInterface
    templateCache map[string]*template.Template
    formDecoder *form.Decoder
    sessionManager *scs.SessionManager
}

func main(){
    //add := flag.String("code", ":4000", "HTTP network address")

    //dsn := flag.String("dsn", "backgo:ppp@/dropbox?parseTime=true", "MySQL data source name")

    flag.Parse()
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = *flag.String("dsn", "postgres://web:pass@localhost:5432/goems?sslmode=disable", "PostgreSQL DSN")
    }
    
    port := os.Getenv("PORT")
if port == "" {
    port = "4000" // local fallback
}
addr := flag.String("addr", ":"+port, "HTTP network address")
    
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
    
    db, err := openDb(dsn)
    if err != nil{
        errLog.Fatal(err)
    }
    
    defer db.Close()
    
    templateCache, err := newTemplateCache("./ui/html/pages/*.tmpl", "./ui/html/base.tmpl","./ui/html/partials/*.tmpl")
    if err != nil{
        errLog.Fatal(err)
    }
    
    formDecoder := form.NewDecoder()
    
    sessionManager := scs.New()
    sessionManager.Store = postgresstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour
    sessionManager.Cookie.Secure = true
    
    app := &application{
        errLog: errLog,
        infoLog: infoLog,
        dbconn : &models.BoxConn{DB: db},
        userconn : &models.UserModel{DB: db},
        templateCache : templateCache,
        formDecoder : formDecoder,
        sessionManager: sessionManager,
    }

    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errLog,
        Handler: app.routes(),
    }
   
    infoLog.Printf("Running on %s",*addr)
    errr := srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
    errLog.Fatal(errr)
    
}
