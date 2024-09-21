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

func TestClientMockCreateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into clients(name, telephone, mail, password) values($1, $2, $3, $4) returning client_id;`).
  		WithArgs("Name", "1234567890", "mail@mail.ru", "123").
  		WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))


	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	client := postgreSQLObjectMother.CreateTestClient()
	err = repository.Create(ctx, client)
	
	assert.NoError(t, err)
		 
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockCreateError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into clients(name, telephone, mail, password) values($1, $2, $3, $4) returning client_id;`).
  		WithArgs("Name", "1234567890", "mail@mail.ru", "123")

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()
	
	client := postgreSQLObjectMother.CreateTestClient()
	err = repository.Create(ctx, client)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where client_id = $1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
		AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_client := postgreSQLObjectMother.CreateTestClient()
	client, err := repository.GetByID(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_client, client)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByIDError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where client_id = $1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByID(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByTelephoneSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where telephone = $1;`).
		WithArgs("1234567890").
		WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
		AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_client := postgreSQLObjectMother.CreateTestClient()
	client, err := repository.GetByTelephone(ctx, "1234567890")

	assert.NoError(t, err)

	assert.Equal(t, new_client, client)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByTelephoneError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where telephone = $1;`).WithArgs("1234567890").WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByTelephone(ctx, "1234567890")

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByTrainingSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where client_id in (select client_id from clients_trainings where training_id=$1);`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
		AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	new_clients := []models.Client{*postgreSQLObjectMother.CreateTestClient()}
	client, err := repository.GetByTraining(ctx, 1)

	assert.NoError(t, err)

	assert.Equal(t, new_clients, client)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockGetByTrainingError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`select * from clients where client_id in (select client_id from clients_trainings where training_id=$1);`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	_, err = repository.GetByTraining(ctx, 1)

	assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockСreateAssignmentSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into clients_trainings(client_id, training_id) values($1, $2) returning client_id;`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.СreateAssignment(ctx, 1, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockСreateAssignmentError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`insert into clients_trainings(client_id, training_id) values($1, $2) returning client_id;`).WithArgs(1, 1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.СreateAssignment(ctx, 1, 1)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockDeleteAssignmentSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`delete from clients_trainings where client_id=$1 and training_id=$2 returning client_id;`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.DeleteAssignment(ctx, 1, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClientMockDeleteAssignmentError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`delete from clients_trainings where client_id=$1 and training_id=$2 returning client_id;`).WithArgs(1, 1)

	dbx := sqlx.NewDb(db, "pgx")
	repository := NewClientPostgreSQLRepository(dbx)
	ctx := context.Background()

	err = repository.DeleteAssignment(ctx, 1, 1)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
