package fetch

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"time"
)

const goodPath = "../InHouseAd/internal/fetch/goods.json"
const categoryPath = "../InHouseAd/internal/fetch/categories.json"

type Fetch interface {
	FetchGoods(ctx context.Context) error
	FetchCategories(ctx context.Context) error
	Init(ctx context.Context, duration time.Duration)
}

type fetch struct {
	gRepo repository.Good
	cRepo repository.Category
	timer *time.Timer
}

func NewFetch(gRepo repository.Good, cRepo repository.Category, timer *time.Timer) Fetch {
	return &fetch{cRepo: cRepo, gRepo: gRepo, timer: timer}
}

func (f *fetch) FetchGoods(ctx context.Context) error {
	absolutePath, err := filepath.Abs(goodPath)
	if err != nil {
		return err
	}

	fileContents, err := os.ReadFile(absolutePath)
	if err != nil {
		return err
	}

	var goodsList []model.Good
	err = json.Unmarshal(fileContents, &goodsList)
	if err != nil {
		return err
	}

	for _, good := range goodsList {
		err = f.gRepo.Create(ctx, &good)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fetch) FetchCategories(ctx context.Context) error {
	absolutePath, err := filepath.Abs(categoryPath)
	if err != nil {
		return err
	}

	fileContents, err := os.ReadFile(absolutePath)
	if err != nil {
		return err
	}

	var categoryList []model.Category
	err = json.Unmarshal(fileContents, &categoryList)
	if err != nil {
		return err
	}

	for _, category := range categoryList {
		err = f.cRepo.Create(ctx, &category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fetch) fetchAll(ctx context.Context) {
	err := f.FetchGoods(ctx)
	if err != nil {
		log.Error().Msgf("unable to fetch goods: %v", err)
		return
	}

	err = f.FetchCategories(ctx)
	if err != nil {
		log.Error().Msgf("unable to fetch categories: %v", err)
		return
	}
}

func (f *fetch) Init(ctx context.Context, duration time.Duration) {
	f.timer = time.AfterFunc(duration, func() {
		f.fetchAll(ctx)
		f.timer.Reset(duration)
	})
}
