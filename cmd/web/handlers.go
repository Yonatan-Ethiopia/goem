package main

import ( "strconv"; "net/http"; "fmt"; _ "html/template"; "github.com/julienschmidt/httprouter"; _ "strings"; _ "unicode/utf8"; "backgo/internal/validator"; "errors"; "backgo/internal/models")


type boxCreateForm struct{
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires_at          int    `form:"expires_at"`
    validator.Validator `form:"-"`
}

type userSignUpForm struct{
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

type userLoginForm struct{
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}


func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request){
    boxes,err := app.dbconn.Latest()
    if err != nil{
        app.serverError(w,err)
    }
    tempData := app.newTemplateData(r)
    tempData.Boxes = boxes
    app.render(w, 200, "home.tmpl", tempData)
}

func (app *application) apiCreatePost( w http.ResponseWriter, r *http.Request){
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7
    
    id, err := app.dbconn.Insert(title, content, expires)
    if err != nil{
        app.infoLog.Printf("Not succefull")
        app.serverError(w, err)
        return
    }
    
    http.Redirect(w,r, fmt.Sprintf("/still/%d", id), http.StatusSeeOther)
    w.Write([]byte("Data inserted"))
}

func (app *application) formCreateGet( w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = &boxCreateForm{
        Expires_at: 365,
    }
    app.render(w, 200, "create.tmpl", data)
}

func (app *application) formCreatePost( w http.ResponseWriter, r *http.Request){

    var form boxCreateForm
    
    err := app.decodePostForm(r, &form)
    if err != nil{
        fmt.Printf("There is an error")
        app.clientError(w, http.StatusBadRequest)
        return
    }
    fmt.Printf("The expiry date is %d", form.Expires_at)
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer then 100 characters")
    form.CheckField(validator.PermittedInt(form.Expires_at, 1,7,365), "expires", "This field cannot be a value other then 1, 7 or 365")
    
    if !form.Valid(){
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }
    
    id, err := app.dbconn.Insert(form.Title, form.Content, form.Expires_at)
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    app.sessionManager.Put(r.Context(),"flash","Box succefully created!")
    http.Redirect(w, r, fmt.Sprintf("/still/%d", id), http.StatusSeeOther)
}

func (app *application) boxViewGet( w http.ResponseWriter, r *http.Request){
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil{
        app.notFound(w)
    }
    
    rec, err := app.dbconn.Get(id)
    if err != nil{
        app.notFound(w)
    }
    
    tempData := app.newTemplateData(r)
    tempData.Box = rec
    app.render(w, http.StatusOK, "view.tmpl", tempData)
}


func (app *application) userSignUpGet(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = userSignUpForm{}
    app.render(w, 200, "signup.tmpl", data)
}
func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request){
    var user userSignUpForm
    
    err := app.decodePostForm(r, &user)
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    user.CheckField(validator.NotBlank(user.Name), "name", "This field cannot be blank")
    user.CheckField(validator.NotBlank(user.Email), "email", "This field cannot be blank")
    user.CheckField(validator.NotBlank(user.Password), "password", "This field cannot be blank")
    
    user.CheckField(validator.MinChars(user.Password, 8), "password", "Password must be atleast 8 characters long")
    user.CheckField(validator.ValidEmailAddress(user.Email, validator.EmailRX), "email", "Invalid email address")
    
    if !user.Valid(){
        fmt.Printf("DEBUG - Decoded user struct: %+v\n", user)
        fmt.Printf("DEBUG - Validation errors:   %+v\n", user.FieldErrors)
         data := app.newTemplateData(r)
         data.Form = user
         app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
         return
    }
    
    err = app.userconn.Insert(user.Name, user.Email, user.Password)
    if err != nil{
        if errors.Is(err, models.ErrDuplicateEmail){
            user.AddFieldError("email", "Email already exists!")
            data := app.newTemplateData(r)
            data.Form = user
            app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
            return
        } else {
            app.serverError(w, err)
        }
        
        return
    }
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login.")
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
    
}
func (app *application) userLogInGet(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = userLoginForm{}
    
    app.render(w, http.StatusOK, "login.tmpl", data)
    
}
func (app *application) userLogInPost(w http.ResponseWriter, r *http.Request){
    
    var user userLoginForm
    var id int
    err := app.decodePostForm(r, &user)
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    
    user.CheckField(validator.NotBlank(user.Email), "email", "This section cannot be empty")
    user.CheckField(validator.ValidEmailAddress(user.Email, validator.EmailRX), "email", "Invalid email address")
    user.CheckField(validator.NotBlank(user.Password), "password", "This section cannot be empty")
    user.CheckField(validator.MinChars(user.Password, 8), "password", "Password needs to be atleast 8 characters long.")
    
    if !user.Valid(){
        data := app.newTemplateData(r)
        data.Form = user
        app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
        return
    }
    
    id, err = app.userconn.Authenticate(user.Email, user.Password)
    if err != nil{
        if errors.Is(err, models.ErrInvalidCredentials){
            user.AddNonFieldError("Incorrect email or password")
            data := app.newTemplateData(r)
            data.Form = user
            app.render(w, http.StatusBadRequest, "login.tmpl", data)
            return
        }
            app.serverError(w, err)
            return
    }
    err = app.sessionManager.RenewToken(r.Context())
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
    
    http.Redirect(w, r, "/create", http.StatusSeeOther)
    
}
func (app *application) userLogOutPost(w http.ResponseWriter, r *http.Request){
    err := app.sessionManager.RenewToken(r.Context())
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    app.sessionManager.Remove(r.Context(), "authenticatedUserId")
    
    app.sessionManager.Put(r.Context(), "flash", "You have successfully logged out! ")
    
    http.Redirect(w,r, "/", http.StatusSeeOther)
}
