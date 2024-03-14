package service

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"fmt"
)

type Category interface {
	Create(ctx context.Context, cat *model.Category) error
	Update(ctx context.Context, cat *model.Category) error
	Get(ctx context.Context) ([]model.Category, error)
	Delete(ctx context.Context, catID uint) error
}

type category struct {
	cRepo repository.Category
	gRepo repository.Good
}

func CategoryService(cRepo repository.Category, gRepo repository.Good) Category {
	return &category{cRepo: cRepo, gRepo: gRepo}
}

func (c *category) Create(ctx context.Context, cat *model.Category) error {
	return c.cRepo.Create(ctx, cat)
}

func (c *category) Update(ctx context.Context, cat *model.Category) error {
	return c.cRepo.Update(ctx, cat)
}

func (c *category) Get(ctx context.Context) ([]model.Category, error) {
	return c.cRepo.Get(ctx)
}

func (c *category) Delete(ctx context.Context, catID uint) error {
	err := c.cRepo.Delete(ctx, catID)
	if err != nil {
		return fmt.Errorf("service.category.Delete.DeleteCategory : %v", err)
	}

	err = c.gRepo.UpdateCategoriesAfterDelete(ctx, catID)
	if err != nil {
		return fmt.Errorf("service.category.Delete.UpdateGoodsCategories : %v", err)
	}

	return nil
}
