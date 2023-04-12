package services

import (
	"errors"
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

	t.Run("Given a new Port at Repository When upserting the Port Then the Port must be created", func(t *testing.T) {
		t.Parallel()

		createWasCalled := false
		updateWasCalled := false
		mockPortRepository := domain.MockPortRepository{
			GetByIDfn: func(id string) (*entities.Port, error) {
				return nil, nil
			},
			Createfn: func(p entities.Port) error {
				createWasCalled = true

				return nil
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

		assert.NoError(t, err, "Error must not be found when upserting an new port")
		assert.True(t, createWasCalled, "Repository's Create method must be called")
		assert.False(t, updateWasCalled, "Repository's Update method must not be called")
	})

	t.Run("Given error is raised When upserting the Port during the query Then an error must be retrieved", func(t *testing.T) {
		t.Parallel()

		createWasCalled := false
		updateWasCalled := false
		mockPortRepository := domain.MockPortRepository{
			GetByIDfn: func(id string) (*entities.Port, error) {
				return nil, errors.New("Error querying Port")
			},
			Createfn: func(p entities.Port) error {
				createWasCalled = true

				return nil
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

		assert.Error(t, err, "Error must be found when upserting an new port")
		assert.False(t, createWasCalled, "Repository's Create method must not be called")
		assert.False(t, updateWasCalled, "Repository's Update method must not be called")
	})

	t.Run("Given error is raised When upserting an existing Port during the update Then an error must be retrieved", func(t *testing.T) {
		t.Parallel()

		createWasCalled := false
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
			Createfn: func(p entities.Port) error {
				createWasCalled = true

				return nil
			},
			Updatefn: func(p entities.Port, filter string) error {
				updateWasCalled = true

				return errors.New("Error during update")
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

		assert.Error(t, err, "Error must be found when upserting an new port")
		assert.False(t, createWasCalled, "Repository's Create method must not be called")
		assert.True(t, updateWasCalled, "Repository's Update method must be called")
	})

	t.Run("Given error is raised When upserting a new Port during the insert Then an error must be retrieved", func(t *testing.T) {
		t.Parallel()

		createWasCalled := false
		updateWasCalled := false
		mockPortRepository := domain.MockPortRepository{
			GetByIDfn: func(id string) (*entities.Port, error) {
				return nil, nil
			},
			Createfn: func(p entities.Port) error {
				createWasCalled = true

				return errors.New("Error during update")
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

		assert.Error(t, err, "Error must be found when upserting an new port")
		assert.True(t, createWasCalled, "Repository's Create method must be called")
		assert.False(t, updateWasCalled, "Repository's Update method must not be called")
	})
}
