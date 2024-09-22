package postgreSQL

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/nkarakotova/lim-core/models"
	"github.com/nkarakotova/lim-core/repositories"

	postgreSQLObjectMother "github.com/nkarakotova/lim-repo/postgreSQL/object_mothers"
	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CoachSuite struct {
	suite.Suite
	db         *sql.DB
	mock       sqlmock.Sqlmock
	repository repositories.CoachRepository
	ctx        context.Context
}

func (s *CoachSuite) BeforeEach(t provider.T) {
	var err error
	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	dbx := sqlx.NewDb(s.db, "pgx")
	s.repository = NewCoachPostgreSQLRepository(dbx)
	s.ctx = context.Background()
}

func (s *CoachSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *CoachSuite) TestCoachMockCreateSuccess(t provider.T) {
	t.Title("CoachMockCreate: Success")
	t.Tags("Coach")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into coaches(name) values($1) returning coach_id;`).
			WithArgs("Name").
			WillReturnRows(sqlmock.NewRows([]string{"coach_id"}).AddRow(1))

		coach := postgreSQLObjectMother.CreateTestCoach()
		err := s.repository.Create(s.ctx, coach)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockCreateFailure(t provider.T) {
	t.Title("CoachMockCreate: Failure")
	t.Tags("Coach")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into coaches(name) values($1) returning coach_id;`).
			WithArgs("Name")

		coach := postgreSQLObjectMother.CreateTestCoach()
		err := s.repository.Create(s.ctx, coach)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetByIDSuccess(t provider.T) {
	t.Title("CoachMockGetByID: Success")
	t.Tags("Coach")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches where coach_id = $1;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
				AddRow(1, "Name"))

		new_coach := postgreSQLObjectMother.CreateTestCoach()
		coach, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_coach, coach)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetByIDFailure(t provider.T) {
	t.Title("CoachMockGetByID: Failure")
	t.Tags("Coach")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches where coach_id = $1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetByNameSuccess(t provider.T) {
	t.Title("CoachMockGetByName: Success")
	t.Tags("Coach")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches where name = $1;`).
			WithArgs("Name").
			WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
				AddRow(1, "Name"))

		new_coach := postgreSQLObjectMother.CreateTestCoach()
		coach, err := s.repository.GetByName(s.ctx, "Name")

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_coach, coach)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetByNameFailure(t provider.T) {
	t.Title("CoachMockGetByName: Failure")
	t.Tags("Coach")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches where name = $1;`).WithArgs("Name").WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByName(s.ctx, "Name")

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetAllSuccess(t provider.T) {
	t.Title("CoachMockGetAll: Success")
	t.Tags("Coach")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches;`).
			WillReturnRows(sqlmock.NewRows([]string{"coach_id", "name"}).
			AddRow(1, "Name"))

		new_coaches := []models.Coach{*postgreSQLObjectMother.CreateTestCoach()}
		coach, err := s.repository.GetAll(s.ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_coaches, coach)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *CoachSuite) TestCoachMockGetAllFailure(t provider.T) {
	t.Title("CoachMockGetAll: Failure")
	t.Tags("Coach")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from coaches;`).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetAll(s.ctx)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCoachSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CoachSuite))
}
