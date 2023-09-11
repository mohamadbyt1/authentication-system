package models

import ("errors")

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserSignup struct {
	Username string `json:"username"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}
type User struct {
	Username string `json:"username"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
	Id int `json:"id"`
}
func NewUser(username string,fname string,lname string,pass string )*User{
	return &User{
		Username: username,
		FirstName: fname,
		LastName: lname,
		Password: pass,
}
}
func (u *UserSignup) ValidateSignup() error{
	if len(u.Username) < 6 || len(u.Password) < 8||len(u.FirstName)< 1||len(u.LastName)< 1{
		return errors.New("not valid")
	}
	return nil
}
func (u *UserLogin) ValidateLogin() error{
	if len(u.Username) < 6 || len(u.Password) < 8{
		return errors.New("not valid")
	}
	return nil
}