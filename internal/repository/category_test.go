package repository_test

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCategoryRepo(t *testing.T) {
	ctx := context.Background()
	repo := repository.CategoryRepo(TestPool)

	cat := model.Category{
		Name: "test",
	}

	err := repo.Create(ctx, &cat)
	assert.NoError(t, err)

	err = repo.Update(ctx, &model.Category{
		ID:   1,
		Name: "testupdate",
	})
	assert.NoError(t, err)

	get, err := repo.Get(ctx)
	assert.NotNil(t, get)
	assert.NoError(t, err)
}
