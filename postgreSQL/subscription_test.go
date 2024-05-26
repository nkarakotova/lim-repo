package postgreSQL

import (
	"context"
	"testing"
	"time"

	"github.com/nkarakotova/lim-core/models"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var testSubscriptionPostgreSQLRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		subscription *models.Subscription
	}
	CheckOutput func(t *testing.T, err error, subscriptionID uint64)
}{
	{
		TestName: "create success test",
		InputData: struct {
			subscription *models.Subscription
		}{&models.Subscription{
			TrainingsNum:          30,
			RemainingTrainingsNum: 30,
			Cost:                  10000,
			StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		}},

		CheckOutput: func(t *testing.T, err error, subscriptionID uint64) {
			assert.NoError(t, err)
			assert.NotEqual(t, subscriptionID, 0)
		},
	},
}

func TestSubscriptionPostgreSQLCreate(t *testing.T) {
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

	for _, tt := range testSubscriptionPostgreSQLRepositoryCreateSuccess {
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

			client, _ := clientRepository.GetByTelephone(ctx, "9262218276")

			subscriptionRepository := CreateSubscriptionPostgreSQLRepository(&fields)
			err := subscriptionRepository.Create(ctx, tt.InputData.subscription, client.ID)

			client, _ = clientRepository.GetByTelephone(ctx, "9262218276")

			tt.CheckOutput(t, err, client.SubscriptionID)
		})
	}
}
