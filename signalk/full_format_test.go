package signalk_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"signalk/signalk"
)

func TestFullFormatShouldMarshalToJsonCorrectlySimple(t *testing.T) {
	vesselID := signalk.VesselIDFromUUID(uuid.New())
	self := signalk.CreateVessel(vesselID)
	full := signalk.FullFormat{
		Version: "1.0",
		Self:    signalk.CreatePath("vessels", vesselID.String()),
		Vessels: signalk.VesselList{
			self,
		},
		Sources: []signalk.Source{},
	}
	expectedJSONString := fmt.Sprintf(
		`{"version":"1.0","self":"vessels.%s","vessels":{"%s":{"uuid":"%s"}},"sources":{}}`,
		self.ID.String(), self.ID.String(), self.ID.String())

	fullBytes, err := json.Marshal(&full)
	require.NoError(t, err)
	fullString := string(fullBytes)

	assert.Equal(t, expectedJSONString, fullString)
}
