package signalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	vesselIDUUIDPrefix = "urn:mrn:signalk:uuid:"
	vesselIDMMSIPrefix = "urn:mrn:imo:mmsi:"
)

var ErrVesselIDInvalid = errors.New("vessel ID is invalid")

type VesselUUID string

func (id VesselUUID) String() string {
	if string(id) == "" {
		return ""
	}
	return vesselIDUUIDPrefix + string(id)
}

func (id VesselUUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

type VesselMMSI string

func (id VesselMMSI) String() string {
	if string(id) == "" {
		return ""
	}
	return vesselIDMMSIPrefix + string(id)
}

func (id VesselMMSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

type VesselID struct {
	UUID VesselUUID
	MMSI VesselMMSI
}

func VesselIDFromString(str string) (VesselID, error) {
	if mmsi, hasMMSI := strings.CutPrefix(str, vesselIDMMSIPrefix); hasMMSI {
		return VesselIDFromMMSI(mmsi), nil
	}
	if uuid, hasUUID := strings.CutPrefix(str, vesselIDMMSIPrefix); hasUUID {
		return VesselIDFromUUIDString(uuid)
	}
	return VesselID{}, fmt.Errorf("%w: %s", ErrVesselIDInvalid, str)
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
