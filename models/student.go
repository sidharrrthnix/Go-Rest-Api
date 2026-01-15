package models

type Student struct {
	ID        int    `json:"id,omitempty" db:"id,omitempty"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Class     string `json:"class" validate:"required"`
}
