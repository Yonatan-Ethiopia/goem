package mocks

import (
        "time"
        "backgo/internal/models"
    )
    
var mockBox = &models.RecRow{
    Id:         1,
    Title:      "An old silent pond",
    Content:    "An old silent pond...",
    Created_at: time.Now(),
    Expires_at: time.Now(),
}

type RecRow struct{}

func (r *RecRow) Insert(title string, content string, expires int) (int, error){
    return 2, nil
}

func (r *RecRow) Get(id int) (*models.RecRow, error){
    switch id{
        case 1:
            return mockBox, nil
        default:
            return nil, models.ErrNoRecord
    }
}

func (r *RecRow) Latest() ([]*models.RecRow, error){
    return []*models.RecRow{mockBox}, nil
}
