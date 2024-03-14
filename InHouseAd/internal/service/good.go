package service

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"sync"
)

type Good interface {
	Create(ctx context.Context, item *model.Good) error
	Update(ctx context.Context, item *model.Good) error
	GetByCategory(ctx context.Context, catID uint) ([]model.Good, error)
	Delete(ctx context.Context, itemID uint) error
}

type good struct {
	gRepo repository.Good

	mtx sync.Mutex
}

func GoodService(gRepo repository.Good) Good {
	return &good{gRepo: gRepo}
}

func (g *good) Create(ctx context.Context, item *model.Good) error {

	return g.gRepo.Create(ctx, item)
}

func (g *good) Update(ctx context.Context, item *model.Good) error {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.gRepo.Update(ctx, item)
}

func (g *good) GetByCategory(ctx context.Context, catID uint) ([]model.Good, error) {
	return g.gRepo.GetByCategory(ctx, catID)
}

func (g *good) Delete(ctx context.Context, itemID uint) error {
	return g.gRepo.Delete(ctx, itemID)
}
