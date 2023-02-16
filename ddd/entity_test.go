package ddd

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	id := "123"
	entity := Create(id, func(inst Entity[string]) Entity[string] {
		return inst
	})

	assert.Equal(t, id, entity.GetId())
	assert.True(t, entity.GetCreatedAt().After(time.Time{}))
	assert.True(t, entity.GetUpdatedAt().After(time.Time{}))
}

func TestEntity_Update(t *testing.T) {
	id := "123"
	entity := Create(id, func(inst Entity[string]) Entity[string] {
		return inst
	})

	createdAt := entity.GetCreatedAt()
	updatedAt := entity.GetUpdatedAt()

	time.Sleep(time.Millisecond)
	entity.Update()

	assert.Equal(t, createdAt, entity.GetCreatedAt())
	assert.True(t, entity.GetUpdatedAt().After(updatedAt))
}
