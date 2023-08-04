package signalk

import (
	"time"

	"signalk/schema"
)

type Service struct {
	root schema.Root
}

func NewService() *Service {
	root := schema.Root{
		Version: "1.0.0",
		Self:    "",
		Vessels: schema.Vessels{},
		Sources: schema.Sources{},
	}

	self := schema.NewVessel()
	root.Vessels.Add(self)

	return &Service{root}
}

func (s *Service) ApplyDelta(delta schema.Delta) error {
	var context string
	if delta.Context == "" {
		context = "vessels." + s.root.Self
	}
	for _, update := range delta.Updates {
		// Process Source
		sourceID := ""

		// Process values
		for _, updateValue := range update.Values {
			path := context + updateValue.Path
			err := s.applyDeltaValue(path, update.Timestamp, sourceID, updateValue.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) applyDeltaValue(path string, timestamp time.Time, source string, value any) error {
	return nil
}
