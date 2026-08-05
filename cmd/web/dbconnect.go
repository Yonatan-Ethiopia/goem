package main

import (
        "database/sql"

)

func openDb(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil{
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
