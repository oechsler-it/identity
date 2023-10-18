package command

import "github.com/oechsler-it/identity/modules/token/domain"

type Revoke struct {
	IdPartial      domain.TokenIdPartial
	RevokingEntity domain.Owner
}
