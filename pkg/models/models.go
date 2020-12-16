package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

type TelemetryWriter interface {
	Add(t Telemetry) (string, error)
}

type TelemetryReader interface {
	Get(id string) (*Telemetry, error)
	GetAll() []Telemetry
}

type HealthChecker interface {
	Alive() (map[string]string, error)
	Ready() (map[string]string, error)
}

type TelemetryReaderWriter interface {
	TelemetryReader
	TelemetryWriter
}

type TelemetryReaderWriterChecker interface {
	TelemetryReaderWriter
	HealthChecker
}

type Position struct {
	// Latitude specifies the north-south position of a location on the globe as degrees
	// and is represented using decimal degrees (DD)
	Latitude float64 `json:"latitude" validate:"required"`

	// Longitude specifies the east-west position of a location on the globe as degrees
	// represented as decimal degrees (DD)
	Longitude float64 `json:"longitude" validate:"required"`

	// Elevation represents the meters above sea level of a object or location in the
	// world.
	Elevation int64 `json:"elevation,omitempty"`
}

type Telemetry struct {
	// Id is the unique object id of the service.  This Id usually represents a physical
	// object in the world such as a vehicle or plane.  This Id should be unique.
	Id string `json:"-"`

	// Position represents the latitude, longitude, and Elevation of a physical object
	// in the world.
	Position Position `json:"position" validate:"required"`

	// Updated stores the last time this object was updated.  This can be useful if you
	// need to expire stale telemetry objects.
	Updated time.Time `json:"updated"`

	// Source is the source of the telemetry point; Reporting sources are required to
	// send their details.
	Source string `json:"source" validate:"required"`

	// ObjectID is a unique id of the telemetry object as it relates to the source.
	ObjectID string `json:"objectId" validate:"required"`

	// Status represents the current status of the object at the time of update
	Status string `json:"status"`
}

func (t *Telemetry) FromJSON(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		return fmt.Errorf("%w: %s", DecodeError, err)
	}

	if err := t.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ValidationError, err)
	}

	return nil
}

func (t *Telemetry) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

func (t *Telemetry) IsExpired() bool {
	return t.Updated.Before(time.Now().Add(-time.Second * 65))
}
