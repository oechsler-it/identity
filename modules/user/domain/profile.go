package domain

type Profile struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}
