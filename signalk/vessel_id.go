package signalk

import (
	"encoding/json"

	"github.com/google/uuid"
)

type VesselUUID string

func (id VesselUUID) String() string {
	if string(id) == "" {
		return ""
	}
	return "urn:mrn:signalk:uuid:" + string(id)
}

func (id VesselUUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

type VesselMMSI string

func (id VesselMMSI) String() string {
	if string(id) == "" {
		return ""
	}
	return "urn:mrn:imo:mmsi:" + string(id)
}

func (id VesselMMSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

type VesselID struct {
	UUID VesselUUID
	MMSI VesselMMSI
}

func VesselIDFromUUIDString(uuidStr string) (VesselID, error) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return VesselID{}, err
	}
	return VesselID{
		UUID: VesselUUID(id.String()),
	}, nil
}

func VesselIDFromUUID(id uuid.UUID) VesselID {
	return VesselID{
		UUID: VesselUUID(id.String()),
	}
}

func VesselIDFromMMSI(mmsi string) VesselID {
	return VesselID{
		MMSI: VesselMMSI(mmsi),
	}
}

func (id VesselID) String() string {
	if id.MMSI != "" {
		return id.MMSI.String()
	}
	return id.UUID.String()
}

func (id VesselID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}
