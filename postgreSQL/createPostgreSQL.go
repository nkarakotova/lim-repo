package postgreSQL

import (
	"database/sql"

	"github.com/nkarakotova/lim-repo/config"
	"github.com/nkarakotova/lim-repo/flags"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/nkarakotova/lim-core/repositories"
)

type PostgresRepositoryFields struct {
	DB     *sql.DB
	Config config.Config
}

func CreatePostgresRepositoryFields(Postgres flags.PostgresFlags, logger *log.Logger) (*PostgresRepositoryFields, error) {
	fields := new(PostgresRepositoryFields)
	var err error
	fields.Config.Postgres = Postgres

	fields.DB, err = fields.Config.Postgres.InitDB(logger)

	if err != nil {
		logger.Error("POSTGRES! Error parse config for postgreSQL")
		return nil, err
	}

	logger.Info("POSTGRES! Successfully create postgres repository fields")

	return fields, nil
}

func CreateClientPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.ClientRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewClientPostgreSQLRepository(dbx)
}

func CreateCoachPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.CoachRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewCoachPostgreSQLRepository(dbx)
}

func CreateDirectionPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.DirectionRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewDirectionPostgreSQLRepository(dbx)
}

func CreateHallPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.HallRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewHallPostgreSQLRepository(dbx)
}

func CreateSubscriptionPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.SubscriptionRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewSubscriptionPostgreSQLRepository(dbx)
}

func CreateTrainingPostgreSQLRepository(fields *PostgresRepositoryFields) repositories.TrainingRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewTrainingPostgreSQLRepository(dbx)
}
