package postgreSQLObjectMother

import (
	"github.com/nkarakotova/lim-core/models"
)

func CreateTestCoach() *models.Coach {
	return &models.Coach{
		ID: 1,
		Name: "Name",
	}
}
