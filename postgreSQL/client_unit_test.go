package postgreSQL

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nkarakotova/lim-core/models"
	"github.com/nkarakotova/lim-core/repositories"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	postgreSQLObjectMother "github.com/nkarakotova/lim-repo/postgreSQL/object_mothers"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ClientSuite struct {
	suite.Suite
	db         *sql.DB
	mock       sqlmock.Sqlmock
	repository repositories.ClientRepository
	ctx        context.Context
}

func (s *ClientSuite) BeforeEach(t provider.T) {
	var err error
	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	dbx := sqlx.NewDb(s.db, "pgx")
	s.repository = NewClientPostgreSQLRepository(dbx)
	s.ctx = context.Background()
}

func (s *ClientSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *ClientSuite) TestClientMockCreateSuccess(t provider.T) {
	t.Title("ClientMockCreate: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into clients(name, telephone, mail, password) values($1, $2, $3, $4) returning client_id;`).
			WithArgs("Name", "1234567890", "mail@mail.ru", "123").
			WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))

		client := postgreSQLObjectMother.CreateTestClient()
		err := s.repository.Create(s.ctx, client)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockCreateFailure(t provider.T) {
	t.Title("ClientMockCreate: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into clients(name, telephone, mail, password) values($1, $2, $3, $4) returning client_id;`).
			WithArgs("Name", "1234567890", "mail@mail.ru", "123")

		client := postgreSQLObjectMother.CreateTestClient()
		err := s.repository.Create(s.ctx, client)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByIDSuccess(t provider.T) {
	t.Title("ClientMockGetByID: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where client_id = $1;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
				AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

		new_client := postgreSQLObjectMother.CreateTestClient()
		client, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_client, client)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByIDFailure(t provider.T) {
	t.Title("ClientMockGetByID: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where client_id = $1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByTelephoneSuccess(t provider.T) {
	t.Title("ClientMockGetByTelephone: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where telephone = $1;`).
			WithArgs("1234567890").
			WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
				AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

		new_client := postgreSQLObjectMother.CreateTestClient()
		client, err := s.repository.GetByTelephone(s.ctx, "1234567890")

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_client, client)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByTelephoneFailure(t provider.T) {
	t.Title("ClientMockGetByTelephone: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where telephone = $1;`).WithArgs("1234567890").WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByTelephone(s.ctx, "1234567890")

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByTrainingSuccess(t provider.T) {
	t.Title("ClientMockGetByTraining: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where client_id in (select client_id from clients_trainings where training_id=$1);`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"client_id", "name", "telephone", "mail", "password"}).
				AddRow(1, "Name", "1234567890", "mail@mail.ru", "123"))

		new_clients := []models.Client{*postgreSQLObjectMother.CreateTestClient()}
		client, err := s.repository.GetByTraining(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_clients, client)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockGetByTrainingFailure(t provider.T) {
	t.Title("ClientMockGetByTraining: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from clients where client_id in (select client_id from clients_trainings where training_id=$1);`).
			WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByTraining(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockCreateAssignmentSuccess(t provider.T) {
	t.Title("ClientMockCreateAssignment: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into clients_trainings(client_id, training_id) values($1, $2) returning client_id;`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))

		err := s.repository.CreateAssignment(s.ctx, 1, 1)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockCreateAssignmentFailure(t provider.T) {
	t.Title("ClientMockCreateAssignment: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into clients_trainings(client_id, training_id) values($1, $2) returning client_id;`).WithArgs(1, 1)

		err := s.repository.CreateAssignment(s.ctx, 1, 1)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockDeleteAssignmentSuccess(t provider.T) {
	t.Title("ClientMockDeleteAssignment: Success")
	t.Tags("Client")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`delete from clients_trainings where client_id=$1 and training_id=$2 returning client_id;`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"client_id"}).AddRow(1))

		err := s.repository.DeleteAssignment(s.ctx, 1, 1)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *ClientSuite) TestClientMockDeleteAssignmentFailure(t provider.T) {
	t.Title("ClientMockDeleteAssignment: Failure")
	t.Tags("Client")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`delete from clients_trainings where client_id=$1 and training_id=$2 returning client_id;`).WithArgs(1, 1)

		err := s.repository.DeleteAssignment(s.ctx, 1, 1)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestClientSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientSuite))
}
