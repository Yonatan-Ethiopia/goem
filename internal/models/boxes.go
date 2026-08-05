package models

import(
    "database/sql"
    "time"
    "errors"
)

//We will use this to hold the data for a new row of a record
type RecRow struct{
    Id int
    Title string
    Content string
    Created_at time.Time
    Expires_at time.Time
}

//We will use this to store the DB model which will hold the connection or allow us to talk to the DB

type BoxConn struct{
    DB *sql.DB
}

type BoxInterface interface{
    Insert(title string, content string, expires_at int) (int, error)
    Get(id int) (*RecRow, error)
    Latest() ([]*RecRow, error)
}

func (c *BoxConn) Insert(title string, content string, expires_at int) (int, error){
    
    stm := "INSERT INTO boxes (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
    
    result, err := c.DB.Exec(stm, title, content, expires_at)
    
    if err != nil{
        return 0, err
    }
    
    id, err := result.LastInsertId()
    
    if err != nil{
        return 0, err
    }
    
    return int(id), nil
}

func (c *BoxConn) Get(id int) (*RecRow, error){
    stm := "SELECT * FROM boxes WHERE UTC_TIMESTAMP < expires AND id = ?"
    
    row := c.DB.QueryRow(stm, id)
    
    rec := &RecRow{ }
    
    err := row.Scan(&rec.Id, &rec.Title, &rec.Content, &rec.Created_at, &rec.Expires_at)
    
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return nil, err
        } else {
            return nil, err
        }
        
    }
    return rec, nil
}

func (c *BoxConn) Latest() ([]*RecRow, error){
    stm := "SELECT * FROM boxes WHERE UTC_TIMESTAMP < expires ORDER BY id DESC LIMIT 10"
    
    rows, err := c.DB.Query(stm)
    
    if err != nil{
        return nil, err
    }
    
    defer rows.Close()
    
    boxArrays := []*RecRow{}
    
    for rows.Next(){
        row := &RecRow{ }
        
        err = rows.Scan(&row.Id, &row.Title, &row.Content, &row.Created_at, &row.Expires_at)
        
        if err != nil{
            return nil, err
        }
        
        boxArrays = append(boxArrays, row)
    }
    
    if err = rows.Err(); err != nil{
        return nil, err
    }
    
    return boxArrays, nil
}
