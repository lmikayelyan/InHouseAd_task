package repository_test

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoodsRepo(t *testing.T) {
	ctx := context.Background()
	repo := repository.GoodRepo(TestPool)

	good := model.Good{
		Name:       "test",
		Categories: []uint{1, 2, 3},
	}

	err := repo.Create(ctx, &good)
	assert.NoError(t, err)

	err = repo.Update(ctx, &model.Good{
		ID:         1,
		Name:       "test",
		Categories: []uint{2, 3},
	})
	assert.NoError(t, err)

	get, err := repo.GetByCategory(ctx, 3)
	assert.NotNil(t, get)
	assert.NoError(t, err)

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}
