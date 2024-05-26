package postgreSQL

import (
	"context"
	"testing"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/models"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var testClientPostgreSQLRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		client *models.Client
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			client *models.Client
		}{&models.Client{Name: "Natali",
			Telephone: "9262218276",
			Mail:      "nka@mail.ru",
			Password:  "123",
			Age:       20,
			Gender:    models.Female}},

		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
}

var testClientPostgreSQLRepositoryCreateFailure = []struct {
	TestName  string
	InputData struct {
		client *models.Client
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "create failure test, telephone number exists",
		InputData: struct {
			client *models.Client
		}{&models.Client{Name: "Natali",
			Telephone: "9262218276",
			Mail:      "nka@mail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

func TestClientPostgreSQLCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	if db == nil {
		return
	}
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	for _, tt := range testClientPostgreSQLRepositoryCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := PostgresRepositoryFields{DB: db}
			ctx := context.Background()

			clientRepository := CreateClientPostgreSQLRepository(&fields)

			err := clientRepository.Create(ctx, tt.InputData.client)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testClientPostgreSQLRepositoryCreateFailure {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := PostgresRepositoryFields{DB: db}
			ctx := context.Background()

			clientRepository := CreateClientPostgreSQLRepository(&fields)

			err := clientRepository.Create(ctx, tt.InputData.client)
			tt.CheckOutput(t, err)
		})
	}
}

var testClientPostgreSQLRepositoryGetByTelephoneSuccess = []struct {
	TestName  string
	InputData struct {
		telephone string
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "get by telephone success test",
		InputData: struct {
			telephone string
		}{telephone: "9262218276"},

		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
}

var testClientPostgreSQLRepositoryGetByTelephoneFailure = []struct {
	TestName  string
	InputData struct {
		telephone string
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "get by telephone failure test, telephone number not exists",
		InputData: struct {
			telephone string
		}{telephone: "9262218275"},

		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, repositoriesErrors.EntityDoesNotExists)
		},
	},
}

func TestClientPostgreSQLGetByTelephone(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	if db == nil {
		return
	}
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	for _, tt := range testClientPostgreSQLRepositoryGetByTelephoneSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := PostgresRepositoryFields{DB: db}
			ctx := context.Background()

			clientRepository := CreateClientPostgreSQLRepository(&fields)
			clientRepository.Create(ctx, &models.Client{Name: "Natali",
				Telephone: "9262218276",
				Mail:      "nka@mail.ru",
				Password:  "123",
				Age:       20,
				Gender:    models.Female})

			_, err := clientRepository.GetByTelephone(ctx, tt.InputData.telephone)
			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testClientPostgreSQLRepositoryGetByTelephoneFailure {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := PostgresRepositoryFields{DB: db}
			ctx := context.Background()

			clientRepository := CreateClientPostgreSQLRepository(&fields)

			_, err := clientRepository.GetByTelephone(ctx, tt.InputData.telephone)
			tt.CheckOutput(t, err)
		})
	}
}
