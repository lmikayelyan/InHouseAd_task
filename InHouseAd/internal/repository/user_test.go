package repository_test

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo(t *testing.T) {
	repo := repository.NewUser(TestPool)

	testUser := model.User{
		UserID:      1,
		Username:    "Test Username",
		Password:    "Test Password",
		PhoneNumber: "+37455884208",
		EMail:       "lmikayelyan@dlink.ru",
	}

	err := repo.Create(context.Background(), testUser)
	assert.NoError(t, err)

	id, err := repo.GetID(context.Background(), testUser.Username)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.EqualValues(t, id, testUser.UserID)

	userHash, err := repo.GetHash(context.Background(), testUser.Username)
	assert.NoError(t, err)
	assert.Len(t, userHash, 13)
	assert.NotNil(t, userHash)
	assert.Equal(t, userHash, testUser.Password)
}
