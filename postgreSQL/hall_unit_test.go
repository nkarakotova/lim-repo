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

func TestHallMockCreateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into halls(number) values($1) returning hall_id;`).
  		WithArgs(1).
  		WillReturnRows(sqlmock.NewRows([]string{"hall_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	hall := postgreSQLObjectMother.CreateTestHall()
	err = repository.Create(ctx, hall)
	
	assert.NoError(t, err)
		 
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockCreateError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into halls(number) values($1) returning hall_id;`).
  		WithArgs(1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	hall := postgreSQLObjectMother.CreateTestHall()
	err = repository.Create(ctx, hall)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls where hall_id=$1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
		AddRow(1, 1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_hall := postgreSQLObjectMother.CreateTestHall()
	hall, err := repository.GetByID(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_hall, hall)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetByIDError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls where hall_id=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByID(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetByNumberSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls where number=$1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
		AddRow(1, 1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_hall := postgreSQLObjectMother.CreateTestHall()
	hall, err := repository.GetByNumber(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_hall, hall)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetByNumberError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls where number=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByNumber(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetAllSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls;`).
		WillReturnRows(sqlmock.NewRows([]string{"hall_id", "number"}).
		AddRow(1, 1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_halls := map[uint64]models.Hall{1: *postgreSQLObjectMother.CreateTestHall()}
	hall, err := repository.GetAll(ctx)

	assert.NoError(t, err)

	assert.Equal(t, new_halls, hall)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHallMockGetAllError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from halls;`).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewHallPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetAll(ctx)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
