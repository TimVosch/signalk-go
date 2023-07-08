package signalk

import (
	"encoding/json"
	"fmt"
)

type NavigationData struct {
	// anchor
	// attitude
	// closestApproach
	// courseGreatCircle
	// courseOverGroundMagnetic
	// courseOverGroundTrue
	// courseRhumbline
	// datetime
	// destination
	// gnss
	// headingCompass
	// headingMagnetic
	// headingTrue
	// leewayAngle
	// lights
	// log
	// magneticDeviation
	// magneticVariation
	// magneticVariationAgeOfService
	// maneuver
	Position Position `json:"position,omitempty"`
	// racing
	// rateOfTurn
	// speedOverGround
	// speedThroughWater
	// speedThroughWaterLongitudinal
	// speedThroughWaterTransverse
	// state
	// trip
}

func (n *NavigationData) applyDeltaUpdateValue(root *Root, path Path, value any) error {
	if path.IsEmpty() {
		return fmt.Errorf("%w: values cannot be merged into navigation data", ErrDeltaUnsupportedUpdate)
	}

	next, _ := path.FirstOut()
	switch next.String() {
	case "position":
		// TODO: Proper checks
		return json.Unmarshal(value.([]byte), &n.Position)
	}

	return fmt.Errorf("%w: for %s in navigation", ErrDeltaPathInvalid, path)
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}
