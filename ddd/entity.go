package ddd

import "time"

type Entity[TId any] struct {
	Id        TId   `json:"id" validate:"required"`
	CreatedAt int64 `json:"created_at" validate:"required,gte=0"`
	UpdatedAt int64 `json:"updated_at" validate:"required,gte=0"`
}

func (e *Entity[TId]) GetId() TId {
	return e.Id
}

func (e *Entity[TId]) GetCreatedAt() time.Time {
	return time.UnixMilli(e.CreatedAt)
}

func (e *Entity[TId]) GetUpdatedAt() time.Time {
	return time.UnixMilli(e.UpdatedAt)
}

func Create[TId any, TInstance any](id TId, ctor func(e Entity[TId]) TInstance) TInstance {
	e := Entity[TId]{
		Id: id,
	}
	e.CreatedAt = time.Now().UnixMilli()
	e.UpdatedAt = e.CreatedAt
	return ctor(e)
}

func (e *Entity[TId]) Update() {
	e.UpdatedAt = time.Now().UnixMilli()
}
