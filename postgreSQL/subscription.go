package postgreSQL

import (
	"context"
	"database/sql"
	"time"

	"github.com/nkarakotova/lim-core/repositories"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"

	"github.com/nkarakotova/lim-core/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

type SubscriptionPostgreSQL struct {
	ID                    uint64    `db:"subscription_id"`
	TrainingsNum          uint64    `db:"trainings_num"`
	RemainingTrainingsNum uint64    `db:"remaining_trainings_num"`
	Cost                  uint64    `db:"cost"`
	StartDate             time.Time `db:"start_date"`
	EndDate               time.Time `db:"end_date"`
}

type SubscriptionPostgreSQLRepository struct {
	db         *sqlx.DB
	txResolver *trmsqlx.CtxGetter
}

func NewSubscriptionPostgreSQLRepository(db *sqlx.DB) repositories.SubscriptionRepository {
	return &SubscriptionPostgreSQLRepository{db: db, txResolver: trmsqlx.DefaultCtxGetter}
}

func (s *SubscriptionPostgreSQLRepository) Create(ctx context.Context, subscription *models.Subscription, clientID uint64) error {
	tr, _ := s.db.Beginx()

	query := `insert into subscriptions(trainings_num, remaining_trainings_num, cost, start_date, end_date) values($1, $2, $3, $4, $5) returning subscription_id;`

	err := s.txResolver.DefaultTrOrDB(ctx, s.db).QueryRowxContext(ctx, query, subscription.TrainingsNum, subscription.RemainingTrainingsNum, subscription.Cost, subscription.StartDate, subscription.EndDate).Scan(&subscription.ID)
	if err != nil {
		tr.Rollback()
		return err
	}

	query = `update clients set subscription_id = $1 where client_id = $2 returning client_id;`
	err = s.txResolver.DefaultTrOrDB(ctx, s.db).QueryRowxContext(ctx, query, subscription.ID, clientID).Scan(&clientID)
	if err != nil {
		tr.Rollback()
		return err
	}

	tr.Commit()

	return nil
}

func (s *SubscriptionPostgreSQLRepository) GetByID(ctx context.Context, id uint64) (*models.Subscription, error) {
	query := `select * from subscriptions where subscription_id=$1;`

	subscriptionDB := &SubscriptionPostgreSQL{}
	err := s.txResolver.DefaultTrOrDB(ctx, s.db).GetContext(ctx, subscriptionDB, query, id)
	if err == sql.ErrNoRows {
		return nil, repositoriesErrors.EntityDoesNotExists
	} else if err != nil {
		return nil, err
	}

	subscriptionModels := &models.Subscription{}
	err = copier.Copy(subscriptionModels, subscriptionDB)
	if err != nil {
		return nil, err
	}

	return subscriptionModels, nil
}

func (s *SubscriptionPostgreSQLRepository) ReduceRemainingTrainingsNum(ctx context.Context, id uint64) error {
	query := `update subscriptions set remaining_trainings_num = remaining_trainings_num - 1 where subscription_id=$1 returning subscription_id;`

	err := s.txResolver.DefaultTrOrDB(ctx, s.db).QueryRowxContext(ctx, query, id).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionPostgreSQLRepository) IncreaseRemainingTrainingsNum(ctx context.Context, id uint64) error {
	query := `update subscriptions set remaining_trainings_num = remaining_trainings_num + 1 where subscription_id=$1 returning subscription_id;`

	err := s.txResolver.DefaultTrOrDB(ctx, s.db).QueryRowxContext(ctx, query, id).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
