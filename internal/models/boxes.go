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
    
    stmt := `INSERT INTO boxes (title, content, created, expires)
VALUES($1, $2, NOW() AT TIME ZONE 'UTC', (NOW() AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 day'))
RETURNING id`
    
    var id int
err := c.DB.QueryRow(stmt, title, content, expires_at).Scan(&id)
if err != nil {
    return 0, err
}
return id, nil
}

func (c *BoxConn) Get(id int) (*RecRow, error){
    stmt := "SELECT * FROM boxes WHERE (NOW() AT TIME ZONE 'UTC') < expires AND id = $1"
    
    row := c.DB.QueryRow(stmt, id)
    
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
    stmt := "SELECT * FROM boxes WHERE (NOW() AT TIME ZONE 'UTC') < expires ORDER BY id DESC LIMIT 10"
    
    rows, err := c.DB.Query(stmt)
    
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
