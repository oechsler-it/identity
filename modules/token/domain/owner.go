package domain

type Owner struct {
	UserId UserId `validate:"required"`
}
