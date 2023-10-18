package domain

type TokenId string

func (i TokenId) GetPartial() TokenIdPartial {
	return NewTokenIdPartial(i)
}
