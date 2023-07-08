package signalk_test

import (
	"encoding/json"
	"signalk/signalk"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullFormatShouldMarshalToJsonCorrectlySimple(t *testing.T) {
	vesselID := signalk.VesselIDFromUUID(uuid.New())
	self := signalk.CreateVessel(vesselID)
	full := signalk.FullFormat{
		Version: "1.0",
		Self:    signalk.CreatePath("vessels", vesselID.String()),
		Vessels: signalk.VesselList{
			vesselID: self,
		},
		Sources: []signalk.Source{},
	}
	expectedJSONString := ""

	fullBytes, err := json.Marshal(&full)
	require.NoError(t, err)
	fullString := string(fullBytes)

	assert.Equal(t, expectedJSONString, fullString)
}
