package repository

import (
	"InHouseAd/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Good interface {
	Create(ctx context.Context, item *model.Good) error
	Update(ctx context.Context, item *model.Good) error
	Delete(ctx context.Context, itemID uint) error
	GetByCategory(ctx context.Context, catID uint) ([]model.Good, error)
	UpdateCategoriesAfterDelete(ctx context.Context, catID uint) error
}

type good struct {
	pool *pgxpool.Pool
}

func GoodRepo(pool *pgxpool.Pool) Good {
	return &good{pool: pool}
}

// Firstly, good object is added to
// the database, then the categories which the goods object belongs to are added
// to 'goods_by_categories' table using batch
func (g *good) Create(ctx context.Context, item *model.Good) error {
	var idList []uint
	for _, cat := range item.Categories {
		err := g.checkCategory(ctx, cat)
		if err != nil {
			return err
		}
		idList = append(idList, cat)
	}

	queryStr := "insert into goods (name) values ($1)"
	_, err := g.pool.Exec(ctx, queryStr, item.Name)
	if err != nil {
		return fmt.Errorf("goods.Create.Exec: %v", err)
	}

	getIDquery := "select(currval('goods_id_seq'))"
	var itemID uint
	err = g.pool.QueryRow(ctx, getIDquery).Scan(&itemID)
	if err != nil {
		return fmt.Errorf("goods.Create.getID.QueryRow: %v", err)
	}

	batchQueryStr := "insert into goods_by_categories(goods_id, category_id) values($1, $2)"
	batch, err := g.prepareBatch(idList, itemID, batchQueryStr)
	if err != nil {
		return fmt.Errorf("goods.Create.prepareBatch: %v", err)
	}

	br := g.pool.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("goods.Create.batchExec: %v", err)
	}

	err = br.Close()
	if err != nil {
		return fmt.Errorf("goods.Create.Batch.Close: %v", err)
	}

	return nil
}

// All the goods entries are updated, and then categories are updated inside
// 'goods_by_categories' table
func (g *good) Update(ctx context.Context, item *model.Good) error {
	tx, err := g.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer func(tx pgx.Tx, ctx context.Context) {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			log.Warn().Msgf("goods.Update.DeferTx : %v", txErr)
		}
	}(tx, ctx)

	var idList []uint
	for _, cat := range item.Categories {
		err := g.checkCategory(ctx, cat)
		if err != nil {
			return err
		}

		idList = append(idList, cat)
	}

	queryStr := "update goods set name=$1 where id=$2"
	_, err = tx.Exec(ctx, queryStr, item.Name, item.ID)
	if err != nil {
		return fmt.Errorf("goods.Update.Exec: %v", err)
	}

	batchQueryStr := "update goods_by_categories set category_id=$1 where goods_id=$2"
	batch, err := g.prepareBatch(idList, item.ID, batchQueryStr)
	if err != nil {
		return fmt.Errorf("goods.Update.prepareBatch: %v", err)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("goods.Update.batchExec: %v", err)
	}

	err = br.Close()
	if err != nil {
		return fmt.Errorf("goods.Update.batch.Close: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("goods.Update.Commit: %v", err)
	}
	return nil
}

// Goods object is deleted, and then the corresponding list of it's
// categories is deleted from 'goods_by_categories' table
func (g *good) Delete(ctx context.Context, itemID uint) error {
	queryStr := "delete from goods where id=$1"
	_, err := g.pool.Exec(ctx, queryStr, itemID)
	if err != nil {
		return fmt.Errorf("goods.delete: %v", err)
	}

	queryStr = "delete from goods_by_categories where goods_id=$1"
	_, err = g.pool.Exec(ctx, queryStr, itemID)
	if err != nil {
		return fmt.Errorf("goodsByCategories.delete: %v", err)
	}

	return err
}

func (g *good) GetByCategory(ctx context.Context, catID uint) ([]model.Good, error) {
	queryStr := `
        select goods.id, goods.name
        from goods_by_categories
        join goods ON goods_by_categories.goods_id = goods.id
        where goods_by_categories.category_id = $1
    `

	rows, err := g.pool.Query(ctx, queryStr, catID)
	if err != nil {
		return nil, fmt.Errorf("goods.GetByCategory.Query: %v", err)
	}
	defer rows.Close()

	var goodsList []model.Good
	for rows.Next() {
		var item model.Good
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, fmt.Errorf("goods.GetByCategory.Rows.Scan: %v", err)
		}
		goodsList = append(goodsList, item)
	}

	return goodsList, nil
}

func (g *good) UpdateCategoriesAfterDelete(ctx context.Context, catID uint) error {
	queryStr := "delete from goods_by_categories where category_id=$1"
	_, err := g.pool.Exec(ctx, queryStr, catID)

	return err
}

func (g *good) checkCategory(ctx context.Context, catID uint) error {
	var check bool
	queryStr := "select exists(select 1 from categories where id=$1)"
	err := g.pool.QueryRow(ctx, queryStr, catID).Scan(&check)

	return err
}

func (g *good) prepareBatch(list []uint, id uint, query string) (*pgx.Batch, error) {
	batch := &pgx.Batch{}
	for _, listId := range list {
		batch.Queue(query, id, listId)
	}

	return batch, nil
}
