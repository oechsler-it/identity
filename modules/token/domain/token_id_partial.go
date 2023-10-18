package domain

type TokenIdPartial struct {
	First string
	Last  string
}

func NewTokenIdPartial(id TokenId) TokenIdPartial {
	return TokenIdPartial{
		First: string(id[:8]),
		Last:  string(id[len(id)-8:]),
	}
}

func NewTokenIdPartialFromString(partial string) (TokenIdPartial, error) {
	if len(partial) == 64 {
		tokenId := TokenId(partial)
		return NewTokenIdPartial(tokenId), nil
	}
	if len(partial) != 16 {
		return TokenIdPartial{}, ErrInvalidTokenIdPartial
	}

	return TokenIdPartial{
		First: partial[len(partial)-8:],
		Last:  partial[:8],
	}, nil
}

func (p TokenIdPartial) String() string {
	return p.Last + p.First
}
