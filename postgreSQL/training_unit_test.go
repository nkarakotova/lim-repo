package postgreSQL

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/models"
	"github.com/nkarakotova/lim-repo/postgreSQL/object_mothers"
	"github.com/stretchr/testify/assert"
)

func TestTrainingMockCreateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into trainings(coach_id, hall_id, name, date_time, places_num) values($1, $2, $3, $4, $5) returning training_id;`).
		WithArgs(1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10).
		WillReturnRows(sqlmock.NewRows([]string{"training_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	training := postgreSQLObjectMother.CreateTestTraining()
	err = repository.Create(ctx, training)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockCreateError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into trainings(coach_id, hall_id, name, date_time, places_num) values($1, $2, $3, $4, $5) returning training_id;`).
		WithArgs(1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	training := postgreSQLObjectMother.CreateTestTraining()
	err = repository.Create(ctx, training)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockDeleteSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`delete from trainings where training_id=$1 returning training_id;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"training_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.Delete(ctx, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`delete from trainings where training_id=$1 returning training_id;`).
		WithArgs(1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.Delete(ctx, 1)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where training_id=$1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
		AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_training := postgreSQLObjectMother.CreateTestTraining()
	training, err := repository.GetByID(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_training, training)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetByIDError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where training_id=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByID(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByClientSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where training_id in (select training_id from clients_trainings where client_id=$1);`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
		AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
	trainings, err := repository.GetAllByClient(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_trainings, trainings)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByClientError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where training_id in (select training_id from clients_trainings where client_id=$1);`).
	WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAllByClient(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByCoachOnDateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where coach_id=$1 and date_time::date=$2::date;`).
		WithArgs(1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
		WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
		AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
	trainings, err := repository.GetAllByCoachOnDate(ctx, 1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

	assert.NoError(t, err)

	assert.Equal(t, new_trainings, trainings)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByCoachOnDateError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where coach_id=$1 and date_time::date=$2::date;`).
	WithArgs(1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAllByCoachOnDate(ctx, 1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByDateTimeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where date_time=$1;`).
		WithArgs(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
		WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
		AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
	trainings, err := repository.GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

	assert.NoError(t, err)

	assert.Equal(t, new_trainings, trainings)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllByDateTimeError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where date_time=$1;`).
	WithArgs(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllBetweenDateTimeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where date_time between $1 and $2;`).
		WithArgs(time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC)).
		WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
		AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
	trainings, err := repository.GetAllBetweenDateTime(ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC))

	assert.NoError(t, err)

	assert.Equal(t, new_trainings, trainings)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockGetAllBetweenDateTimeError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from trainings where date_time between $1 and $2;`).
	WithArgs(time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAllBetweenDateTime(ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC))

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockReduceAvailablePlacesNumSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`update trainings set available_places_num = available_places_num - 1 where training_id=$1 returning training_id;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"training_id"}).
		AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.ReduceAvailablePlacesNum(ctx, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockReduceAvailablePlacesNumError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`update trainings set available_places_num = available_places_num - 1 where training_id=$1 returning training_id;`).WithArgs(1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.ReduceAvailablePlacesNum(ctx, 1)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockIncreaseAvailablePlacesNumSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`update trainings set available_places_num = available_places_num + 1 where training_id=$1 returning training_id;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"training_id"}).
		AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.IncreaseAvailablePlacesNum(ctx, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrainingMockIncreaseAvailablePlacesNumError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`update trainings set available_places_num = available_places_num + 1 where training_id=$1 returning training_id;`).WithArgs(1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewTrainingPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.IncreaseAvailablePlacesNum(ctx, 1)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
