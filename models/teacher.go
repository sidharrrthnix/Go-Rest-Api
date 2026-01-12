package models

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Class     string `json:"class"`
	Subject   string `json:"subject"`
}
