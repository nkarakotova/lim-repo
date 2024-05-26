package postgreSQL

import (
	"context"
	"database/sql"

	"github.com/nkarakotova/lim-core/repositories"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"

	"github.com/nkarakotova/lim-core/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

type HallPostgreSQL struct {
	ID       uint64 `db:"hall_id"`
	Number   uint64 `db:"number"`
	Capacity uint64 `db:"capacity"`
}

type HallPostgreSQLRepository struct {
	db         *sqlx.DB
	txResolver *trmsqlx.CtxGetter
}

func NewHallPostgreSQLRepository(db *sqlx.DB) repositories.HallRepository {
	return &HallPostgreSQLRepository{db: db, txResolver: trmsqlx.DefaultCtxGetter}
}

func (h *HallPostgreSQLRepository) Create(ctx context.Context, hall *models.Hall) error {
	query := `insert into halls(number, capacity) values($1, $2) returning hall_id;`

	err := h.txResolver.DefaultTrOrDB(ctx, h.db).QueryRowxContext(ctx, query, hall.Number, hall.Capacity).Scan(&hall.ID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HallPostgreSQLRepository) GetByID(ctx context.Context, id uint64) (*models.Hall, error) {
	query := `select * from halls where hall_id=$1;`

	hallDB := &HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).GetContext(ctx, hallDB, query, id)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := &models.Hall{}
	err = copier.Copy(hallModels, hallDB)
	if err != nil {
		return nil, err
	}

	return hallModels, nil
}

func (h *HallPostgreSQLRepository) GetByNumber(ctx context.Context, number uint64) (*models.Hall, error) {
	query := `select * from halls where number=$1;`

	hallDB := &HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).GetContext(ctx, hallDB, query, number)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := &models.Hall{}
	err = copier.Copy(hallModels, hallDB)
	if err != nil {
		return nil, err
	}

	return hallModels, nil
}

func (h *HallPostgreSQLRepository) GetAll(ctx context.Context) (map[uint64]models.Hall, error) {
	query := `select * from halls;`

	hallDB := []HallPostgreSQL{}
	err := h.txResolver.DefaultTrOrDB(ctx, h.db).SelectContext(ctx, &hallDB, query)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	hallModels := make(map[uint64]models.Hall)
	for i := range hallDB {
		hall := models.Hall{}
		err = copier.Copy(&hall, &hallDB[i])
		if err != nil {
			return nil, err
		}

		hallModels[hall.ID] = hall
	}

	return hallModels, nil
}
