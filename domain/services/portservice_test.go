package services

import (
	"testing"

	"github.com/cassiuspaim/portimporter/domain"
	"github.com/cassiuspaim/portimporter/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestUpsertPort(t *testing.T) {
	t.Parallel()

	t.Run("Given an existing Port at Repository When upserting the Port Then the Port must be updated", func(t *testing.T) {
		t.Parallel()

		updateWasCalled := false
		mockPortRepository := domain.MockPortRepository{
			GetByIDfn: func(id string) (*entities.Port, error) {
				port := entities.NewPort(
					"id",
					"name",
					"city",
					"country",
					[]string{"alias1", "alias2"},
					[]string{"region1", "region2"},
					[]float64{43.434343434, 35.2423434},
					"province",
					"timezone",
					[]string{"unloc1", "unloc2"},
					"code")

				return &port, nil
			},

			Updatefn: func(p entities.Port, filter string) error {
				updateWasCalled = true

				return nil
			},
		}

		portService := NewPortService(mockPortRepository)
		err := portService.Upsert(entities.NewPort(
			"id",
			"name",
			"city",
			"country",
			[]string{"alias1", "alias2"},
			[]string{"region1", "region2"},
			[]float64{43.434343434, 35.2423434},
			"province",
			"timezone",
			[]string{"unloc1", "unloc2"},
			"code"))

		assert.NoError(t, err, "Error must not be found when upserting an existing port")
		assert.True(t, updateWasCalled, "Repository's Update method must be called")
	})
}
