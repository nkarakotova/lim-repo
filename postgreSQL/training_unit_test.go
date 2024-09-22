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
	"github.com/nkarakotova/lim-core/repositories"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type TrainingSuite struct {
	suite.Suite
	db         *sql.DB
	mock       sqlmock.Sqlmock
	repository repositories.TrainingRepository
	ctx        context.Context
}

func (s *TrainingSuite) BeforeEach(t provider.T) {
	var err error
	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	dbx := sqlx.NewDb(s.db, "pgx")
	s.repository = NewTrainingPostgreSQLRepository(dbx)
	s.ctx = context.Background()
}

func (s *TrainingSuite) AfterEach(t provider.T) {
	s.db.Close()
}

func (s *TrainingSuite) TestTrainingMockCreateSuccess(t provider.T) {
	t.Title("TrainingMockCreate: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into trainings(coach_id, hall_id, name, date_time, places_num) values($1, $2, $3, $4, $5) returning training_id;`).
			WithArgs(1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10).
			WillReturnRows(sqlmock.NewRows([]string{"training_id"}).AddRow(1))

		training := postgreSQLObjectMother.CreateTestTraining()
		err := s.repository.Create(s.ctx, training)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockCreateFailure(t provider.T) {
	t.Title("TrainingMockCreate: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`insert into trainings(coach_id, hall_id, name, date_time, places_num) values($1, $2, $3, $4, $5) returning training_id;`).
			WithArgs(1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10)	


		training := postgreSQLObjectMother.CreateTestTraining()
		err := s.repository.Create(s.ctx, training)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockDeleteSuccess(t provider.T) {
	t.Title("TrainingMockDelete: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`delete from trainings where training_id=$1 returning training_id;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"training_id"}).AddRow(1))

		err := s.repository.Delete(s.ctx, 1)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockDeleteFailure(t provider.T) {
	t.Title("TrainingMockDelete: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`delete from trainings where training_id=$1 returning training_id;`).
			WithArgs(1)

		err := s.repository.Delete(s.ctx, 1)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetByIDSuccess(t provider.T) {
	t.Title("TrainingMockGetByID: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where training_id=$1;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
			AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))

		new_training := postgreSQLObjectMother.CreateTestTraining()
		training, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_training, training)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetByIDFailure(t provider.T) {
	t.Title("TrainingMockGetByID: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where training_id=$1;`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetByID(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByClientSuccess(t provider.T) {
	t.Title("TrainingMockGetAllByClient: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where training_id in (select training_id from clients_trainings where client_id=$1);`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
			AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))


		new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
		trainings, err := s.repository.GetAllByClient(s.ctx, 1)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_trainings, trainings)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByClientFailure(t provider.T) {
	t.Title("TrainingMockGetAllByClient: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where training_id in (select training_id from clients_trainings where client_id=$1);`).
			WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetAllByClient(s.ctx, 1)

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByCoachOnDateSuccess(t provider.T) {
	t.Title("TrainingMockGetAllByCoachOnDate: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where coach_id=$1 and date_time::date=$2::date;`).
			WithArgs(1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
			WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
			AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))


		new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
		trainings, err := s.repository.GetAllByCoachOnDate(s.ctx, 1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_trainings, trainings)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByCoachOnDateFailure(t provider.T) {
	t.Title("TrainingMockGetAllByCoachOnDate: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where coach_id=$1 and date_time::date=$2::date;`).
			WithArgs(1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)

		_, err := s.repository.GetAllByCoachOnDate(s.ctx, 1, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByDateTimeSuccess(t provider.T) {
	t.Title("TrainingMockGetAllByDateTime: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where date_time=$1;`).
			WithArgs(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
			WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
			AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))


		new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
		trainings, err := s.repository.GetAllByDateTime(s.ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_trainings, trainings)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllByDateTimeFailure(t provider.T) {
	t.Title("TrainingMockGetAllByDateTime: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where date_time=$1;`).
			WithArgs(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)
		
		_, err := s.repository.GetAllByDateTime(s.ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllBetweenDateTimeSuccess(t provider.T) {
	t.Title("TrainingMockGetAllBetweenDateTime: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where date_time between $1 and $2;`).
			WithArgs(time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC)).
			WillReturnRows(sqlmock.NewRows([]string{"training_id", "coach_id", "hall_id", "name", "date_time", "places_num"}).
			AddRow(1, 1, 1, "Name", time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC), 10))


		new_trainings := []models.Training{*postgreSQLObjectMother.CreateTestTraining()}
		trainings, err := s.repository.GetAllBetweenDateTime(s.ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(new_trainings, trainings)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockGetAllBetweenDateTimeFailure(t provider.T) {
	t.Title("TrainingMockGetAllBetweenDateTime: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`select * from trainings where date_time between $1 and $2;`).
			WithArgs(time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC)).WillReturnError(sql.ErrNoRows)
		
		_, err := s.repository.GetAllBetweenDateTime(s.ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 10, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().ErrorIs(err, repositoriesErrors.EntityDoesNotExists)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockReduceAvailablePlacesNumSuccess(t provider.T) {
	t.Title("TrainingMockReduceAvailablePlacesNum: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`update trainings set available_places_num = available_places_num - 1 where training_id=$1 returning training_id;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"training_id"}).
			AddRow(1))

		err := s.repository.ReduceAvailablePlacesNum(s.ctx, 1)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockReduceAvailablePlacesNumFailure(t provider.T) {
	t.Title("TrainingMockReduceAvailablePlacesNum: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`update trainings set available_places_num = available_places_num - 1 where training_id=$1 returning training_id;`).WithArgs(1)
		
		err := s.repository.ReduceAvailablePlacesNum(s.ctx, 1)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockIncreaseAvailablePlacesNumSuccess(t provider.T) {
	t.Title("TrainingMockIncreaseAvailablePlacesNum: Success")
	t.Tags("Training")
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`update trainings set available_places_num = available_places_num + 1 where training_id=$1 returning training_id;`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"training_id"}).
			AddRow(1))

		err := s.repository.IncreaseAvailablePlacesNum(s.ctx, 1)

		sCtx.Assert().NoError(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *TrainingSuite) TestTrainingMockIncreaseAvailablePlacesNumFailure(t provider.T) {
	t.Title("TrainingMockIncreaseAvailablePlacesNum: Failure")
	t.Tags("Training")
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		s.mock.ExpectQuery(`update trainings set available_places_num = available_places_num + 1 where training_id=$1 returning training_id;`).WithArgs(1)
		
		err := s.repository.IncreaseAvailablePlacesNum(s.ctx, 1)

		sCtx.Assert().Error(err)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTrainingSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainingSuite))
}
