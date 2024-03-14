package repository

import (
	"InHouseAd/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Category interface {
	Create(ctx context.Context, cat *model.Category) error
	Update(ctx context.Context, cat *model.Category) error
	Get(ctx context.Context) ([]model.Category, error)
	Delete(ctx context.Context, catID uint) error
}

type category struct {
	pool *pgxpool.Pool
}

func CategoryRepo(pool *pgxpool.Pool) Category {
	return &category{pool: pool}
}

func (c *category) Create(ctx context.Context, cat *model.Category) error {
	queryStr := "insert into categories (name) VALUES($1)"
	_, err := c.pool.Exec(ctx, queryStr, cat.Name)

	return err
}

func (c *category) Update(ctx context.Context, cat *model.Category) error {
	queryStr := "update categories set name=$1 where id=$2"
	_, err := c.pool.Exec(ctx, queryStr, cat.Name, cat.ID)

	return err
}

func (c *category) Get(ctx context.Context) ([]model.Category, error) {
	queryStr := "select * from categories"
	rows, err := c.pool.Query(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryList []model.Category
	for rows.Next() {
		var item model.Category
		err = rows.Scan(&item.ID,
			&item.Name)

		if err != nil {
			return nil, fmt.Errorf("error occured while scanning 'categories' table rows: %v", err)
		}

		categoryList = append(categoryList, item)
	}

	return categoryList, nil
}

func (c *category) Delete(ctx context.Context, catID uint) error {
	checkQuery := "select goods_id from goods_by_categories where category_id = $1 limit 2"
	rows, err := c.pool.Query(ctx, checkQuery, catID)
	if err != nil {
		return fmt.Errorf("category.Delete.CheckQuery: %v", err)
	}
	defer rows.Close()

	var goodsIDs []uint
	for rows.Next() {
		var goodsID uint
		if err := rows.Scan(&goodsID); err != nil {
			return fmt.Errorf("category.Delete.Rows.Scan: %v", err)
		}
		goodsIDs = append(goodsIDs, goodsID)
	}

	if len(goodsIDs) == 1 {
		// The category is the only category for a single good, so don't delete it
		return fmt.Errorf("category.Delete: Cannot delete the only category for a good")
	}

	queryStr := "delete from categories where id=$1"
	_, err = c.pool.Exec(ctx, queryStr, catID)

	return err
}
