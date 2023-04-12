package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPort(t *testing.T) {
	t.Parallel()
	t.Run("Given parameters When instantiating a new Port Then the properties of the Port should be equal to the parameters", func(t *testing.T) {
		t.Parallel()

		expetedID := "id"
		expetedName := "name"
		expectedCity := "city"
		expectedCountry := "country"
		expectedAlias := []string{"alias1", "alias2"}
		expectedRegion := []string{"region1", "region2"}
		expectedCoordinates := []float64{43.434343434, 35.2423434}
		expectedProvince := "province"
		expectedTimezone := "timezone"
		expectedUnlocs := []string{"unloc1", "unloc2"}
		expectedCode := "code"
		port := NewPort(
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
		assert.Equal(t, expetedID, port.ID, "IDs must be equal")
		assert.Equal(t, expetedName, port.Name, "Names must be equal")
		assert.Equal(t, expectedCity, port.City, "Cities must be equal")
		assert.Equal(t, expectedCountry, port.Country, "Countries must be equal")
		assert.Equal(t, expectedAlias, port.Alias, "Aliases must be equal")
		assert.Equal(t, expectedRegion, port.Regions, "Regions must be equal")
		assert.Equal(t, expectedCoordinates, port.Coordinates, "Coordinates must be equal")
		assert.Equal(t, expectedProvince, port.Province, "Provinces must be equal")
		assert.Equal(t, expectedTimezone, port.Timezone, "Timezones must be equal")
		assert.Equal(t, expectedUnlocs, port.Unlocs, "Unlocs must be equal")
		assert.Equal(t, expectedCode, port.Code, "Codes must be equal")
	})
}
