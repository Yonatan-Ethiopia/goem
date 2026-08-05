package mocks

import (
    "backgo/internal/models"
    )
    
type User struct{}

func (u *User) Insert(name, email, password string) error {
    switch email {
        case "dupe@example.com":
            return models.ErrDuplicateEmail
        default:
            return nil
    }
}

func (u *User) Authenticate(email, password string) (int, error){
    if email == "dupe@exmaple.com" && password == "validPassword"{
         return 1, nil
    }
    
    return 0, models.ErrInvalidCredentials
}

func (u *User) Exists(id int) (bool, error){
    switch id{
        case 1:
            return true, nil
        default: 
            return false, nil
    }
}
