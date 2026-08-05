package main

import ("backgo/internal/models"; "html/template"; "path/filepath"; "time";)

type templateData struct{
    Box *models.RecRow
    Boxes []*models.RecRow
    CurrentYear int
    Form any
    Flash string
    IsAuthenticated bool
    CSRFToken string
}

func humanDate(t time.Time) string {
    if t.IsZero(){
        return ""
    }
    
    return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache(fileGlob, basePath, partialPath string) (map[string]*template.Template, error){
    cache := map[string]*template.Template{}
    
    pages, err := filepath.Glob(fileGlob)
    if err != nil{
        return nil, err
    }
    
    for _,page := range pages {
        name := filepath.Base(page)
    
        ts, err := template.New(name).Funcs(functions).ParseFiles(basePath)
        if err != nil{
            return nil, err
        }
        
        ts, err = ts.ParseGlob(partialPath)
        if err != nil{
            return nil, err
        }
        
        ts, err = ts.ParseFiles(page)
        if err != nil{
            return nil, err
        }
        
        
        cache[name] = ts
    }
    return cache, nil
}
