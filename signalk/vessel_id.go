package signalk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type VesselIDType uint8

const (
	VesselMMSI VesselIDType = iota
	VesselUUID
	VesselURL

	// Prefixes
	VesselIDMMSIPrefix string = "urn:mrn:imo:mmsi:"
	VesselIDUUIDPrefix string = "urn:mrn:signalk:uuid:"
)

var ErrVesselIDInvalid = errors.New("vessel ID is invalid")

func (t VesselIDType) String() string {
	switch t {
	case VesselMMSI:
		return "mmsi"
	case VesselUUID:
		return "uuid"
	case VesselURL:
		return "url"
	}
	return ""
}

type VesselID struct {
	typ  VesselIDType
	mmsi string
	uuid uuid.UUID
	url  string
}

func (id VesselID) String() string {
	switch id.typ {
	case VesselMMSI:
		return VesselIDMMSIPrefix + id.mmsi
	case VesselUUID:
		return VesselIDUUIDPrefix + id.uuid.String()
	case VesselURL:
		return id.url
	}
	return ""
}

func CreateVesselUUID(uuid uuid.UUID) VesselID {
	return VesselID{
		typ:  VesselUUID,
		uuid: uuid,
	}
}

func ParseVesselID(idStr string) (VesselID, error) {
	var id VesselID
	if mmsi, hasMMSI := strings.CutPrefix(idStr, VesselIDMMSIPrefix); hasMMSI {
		id.typ = VesselMMSI
		id.mmsi = mmsi
		return id, nil
	}
	if uuidStr, hasUUID := strings.CutPrefix(idStr, VesselIDUUIDPrefix); hasUUID {
		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			return id, fmt.Errorf("%w: vessel uuid ('%s') is invalid", ErrVesselIDInvalid, uuidStr)
		}
		id.typ = VesselUUID
		id.uuid = uuid
		return id, nil
	}
	return id, fmt.Errorf("%w: id '%s' is not recognized", ErrVesselIDInvalid, idStr)
}
