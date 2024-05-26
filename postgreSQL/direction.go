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

type DirectionPostgreSQL struct {
	ID               uint64        `db:"direction_id"`
	Name             string        `db:"name"`
	Description      string        `db:"description"`
	AcceptableGender models.Gender `db:"acceptable_gender"`
}

type DirectionPostgreSQLRepository struct {
	db         *sqlx.DB
	txResolver *trmsqlx.CtxGetter
}

func NewDirectionPostgreSQLRepository(db *sqlx.DB) repositories.DirectionRepository {
	return &DirectionPostgreSQLRepository{db: db, txResolver: trmsqlx.DefaultCtxGetter}
}

func (d *DirectionPostgreSQLRepository) Create(ctx context.Context, direction *models.Direction) error {
	query := `insert into directions(name, description, acceptable_gender) values($1, $2, $3) returning direction_id;`

	err := d.txResolver.DefaultTrOrDB(ctx, d.db).QueryRowxContext(ctx, query, direction.Name, direction.Description, direction.AcceptableGender).Scan(&direction.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DirectionPostgreSQLRepository) GetByID(ctx context.Context, id uint64) (*models.Direction, error) {
	query := `select * from directions where direction_id=$1;`

	directionDB := &DirectionPostgreSQL{}
	err := d.txResolver.DefaultTrOrDB(ctx, d.db).GetContext(ctx, directionDB, query, id)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	directionModels := &models.Direction{}
	err = copier.Copy(directionModels, directionDB)
	if err != nil {
		return nil, err
	}

	return directionModels, nil
}

func (d *DirectionPostgreSQLRepository) GetByName(ctx context.Context, name string) (*models.Direction, error) {
	query := `select * from directions where name=$1;`

	directionDB := &DirectionPostgreSQL{}
	err := d.txResolver.DefaultTrOrDB(ctx, d.db).GetContext(ctx, directionDB, query, name)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	directionModels := &models.Direction{}
	err = copier.Copy(directionModels, directionDB)
	if err != nil {
		return nil, err
	}

	return directionModels, nil
}

func (d *DirectionPostgreSQLRepository) GetAll(ctx context.Context) ([]models.Direction, error) {
	query := `select * from directions;`

	directionDB := []DirectionPostgreSQL{}
	err := d.txResolver.DefaultTrOrDB(ctx, d.db).SelectContext(ctx, &directionDB, query)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	directionModels := []models.Direction{}
	for i := range directionDB {
		direction := models.Direction{}
		err = copier.Copy(&direction, &directionDB[i])
		if err != nil {
			return nil, err
		}

		directionModels = append(directionModels, direction)
	}

	return directionModels, nil
}
