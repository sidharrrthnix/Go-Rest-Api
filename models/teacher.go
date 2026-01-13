package models

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Class     string `json:"class" validate:"required"`
	Subject   string `json:"subject" validate:"required"`
}
