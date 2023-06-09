package jsonstream

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartWithValidContentFile(t *testing.T) {
	t.Parallel()

	t.Run("Given a content with valid Port object When reading the file Then no error is expected", func(t *testing.T) {
		t.Parallel()

		expectedKey := "AEAJM"
		expectedName := "Name"
		expectedCity := "City"
		expectedCountry := "country"
		expectedAlias := []string{"Alias1", "Alias2"}
		expectedRegions := []string{"Region1", "Region2"}
		expectedCoordinates := []float64{34.434343, 67.354545}
		expectedProvince := "province"
		expectedTimezone := "Asia/Dubai"
		expectedUnlocs := []string{"unloc1"}
		expectedCode := "code"
		portJSON := getJSONPort(expectedName, expectedCity, expectedCountry, expectedAlias, expectedRegions,
			expectedCoordinates, expectedProvince, expectedTimezone, expectedUnlocs, expectedCode)

		fileContent := strings.NewReader(fmt.Sprintf(`{ "%s": %s}`, expectedKey, portJSON))
		stream := NewPortStream()
		go func() {
			stream.Start(fileContent)
		}()

		for entry := range stream.Watch() {
			if entry.Error == nil {
				assert.Equal(t, expectedKey, entry.Key, "Key must be equal")
				assert.Equal(t, expectedName, entry.Data.Name, "Name must be equal")
				assert.Equal(t, expectedCoordinates, entry.Data.Coordinates, "Coordinates must be equal")
				assert.Equal(t, expectedCity, entry.Data.City, "City must be equal")
				assert.Equal(t, expectedProvince, entry.Data.Province, "Province must be equal")
				assert.Equal(t, expectedCountry, entry.Data.Country, "Country must be equal")
				assert.Equal(t, expectedAlias, entry.Data.Alias, "Alias must be equal")
				assert.Equal(t, expectedRegions, entry.Data.Regions, "Regions must be equal")
				assert.Equal(t, expectedUnlocs, entry.Data.Unlocs, "Unlocs must be equal")
				assert.Equal(t, expectedTimezone, entry.Data.Timezone, "Timezone must be equal")
				assert.Equal(t, expectedCode, entry.Data.Code, "Code must be equal")
			}
		}
	})
}

func TestStartWithInvalidContentFile(t *testing.T) {
	t.Parallel()

	name := "Name"
	city := "City"
	country := "country"
	alias := []string{"Alias1", "Alias2"}
	regions := []string{"Region1", "Region2"}
	coordinates := []float64{34.434343, 67.354545}
	province := "province"
	timezone := "Asia/Dubai"
	unlocs := []string{"unloc1"}
	code := "code"
	portJSON := getJSONPort(name, city, country, alias, regions, coordinates, province, timezone, unlocs, code)

	tests := []struct {
		name        string
		fileContent string
	}{
		{
			name:        "Given a content with missing the key When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`{ : %s}`, ""),
		},
		{
			name:        "Given a content with a key without quotation marks When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`{ key : %s }`, ""),
		},
		{
			name:        "Given a content without opening delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`"key" : %s }`, ""),
		},
		{
			name:        "Given a content with unexpected opening delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`[ "key" : %s }`, ""),
		},
		{
			name:        "Given a content with invalid opening delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`| "key" : %s }`, ""),
		},
		{
			name:        "Given a content without closing delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`{ "key" : %s`, portJSON),
		},
		{
			name:        "Given a content with unexpected closing delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`{ "key" : %s ]`, portJSON),
		},
		{
			name:        "Given a content with invalid closing delimiter When reading the file Then an error is expected",
			fileContent: fmt.Sprintf(`{ "key" : %s |`, portJSON),
		},
		{
			name:        "Given a content with an invalid Port When reading the file Then an error is expected",
			fileContent: `{ "key" : x }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stream := NewPortStream()
			go func() {
				stream.Start(strings.NewReader(tt.fileContent))
			}()
			var errorFound error
			for entry := range stream.Watch() {
				if entry.Error != nil {
					errorFound = entry.Error
				}
			}
			assert.Error(t, errorFound)
		})
	}
}

func getJSONPort(name string, city string, country string, alias []string, regions []string,
	coordinates []float64, province string, timezone string, unlocs []string, code string) string {

	jsonAlias := sliceToJSON(alias)
	jsonRegions := sliceToJSON(regions)
	jsonCoordinates := sliceToJSON(coordinates)
	jsonUnlocs := sliceToJSON(unlocs)

	content := fmt.Sprintf(`
{
    "name": "%s",
    "city": "%s",
    "country": "%s",
    "alias": %s,
    "regions": %s,
    "coordinates": %s,
    "province": "%s",
    "timezone": "%s",
    "unlocs": %s,
    "code": "%s"
}`, name, city, country, jsonAlias, jsonRegions,
		jsonCoordinates, province, timezone, jsonUnlocs, code)

	return fmt.Sprintln(content)
}

func sliceToJSON(value interface{}) []byte {
	json, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Error marshalling the value %v. Error: %v", value, err)
	}

	return json
}
