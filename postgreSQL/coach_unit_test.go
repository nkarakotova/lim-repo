package postgreSQL

import (
	"context"
	"testing"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nkarakotova/lim-core/models"

	"github.com/nkarakotova/lim-repo/postgreSQL/object_mothers"
	"github.com/stretchr/testify/assert"
	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
)

func TestCoachMockCreateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into coaches(name) values($1) returning coach_id;`).
  		WithArgs("Name").
  		WillReturnRows(sqlmock.NewRows([]string{"coach_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	coach := postgreSQLObjectMother.CreateTestCoach()
	err = repository.Create(ctx, coach)
	
	assert.NoError(t, err)
		 
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockCreateError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into coaches(name) values($1) returning coach_id;`).
  		WithArgs("Name")

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	coach := postgreSQLObjectMother.CreateTestCoach()
	err = repository.Create(ctx, coach)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches where coach_id = $1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
		AddRow(1, "Name"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_coach := postgreSQLObjectMother.CreateTestCoach()
	coach, err := repository.GetByID(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_coach, coach)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetByIDError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches where coach_id = $1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByID(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetByNameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches where name = $1;`).
		WithArgs("Name").
		WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
		AddRow(1, "Name"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_coach := postgreSQLObjectMother.CreateTestCoach()
	coach, err := repository.GetByName(ctx, "Name")

	assert.NoError(t, err)

	assert.Equal(t, new_coach, coach)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetByNameError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches where name = $1;`).WithArgs("Name").WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByName(ctx, "Name")

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetAllSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches;`).
		WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
		AddRow(1, "Name"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_coaches := []models.Coach{*postgreSQLObjectMother.CreateTestCoach()}
	coach, err := repository.GetAll(ctx)

	assert.NoError(t, err)

	assert.Equal(t, new_coaches, coach)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCoachMockGetAllError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from coaches;`).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewCoachPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAll(ctx)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
