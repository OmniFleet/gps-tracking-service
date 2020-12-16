package inmem

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

func TestAdd(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)
	item := models.Telemetry{
		Id: uuid.New().String(),
		Position: models.Position{
			Latitude:  0.00,
			Longitude: 0.00,
			Elevation: 0,
		},
	}
	id, err := db.Add(item)
	if err != nil {
		t.Errorf("could not add data to db: %s", err.Error())
	}
	if id != item.Id {
		t.Errorf("item id returned is not our id: %s / %s", id, item.Id)
	}
}

func TestGet(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)
	itemIn := models.Telemetry{
		Id: uuid.New().String(),
		Position: models.Position{
			Latitude:  0.00,
			Longitude: 0.00,
			Elevation: 0,
		},
	}
	id, err := db.Add(itemIn)
	if err != nil {
		t.Errorf("error adding to the database: %s", err.Error())
	}
	fmt.Printf("got %s id back\n", id)

	_, err = db.Get(itemIn.Id)
	if err != nil {
		t.Errorf("error getting item out of db: %s", err.Error())
	}
}

func TestGetNotFound(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)

	result, err := db.Get("foobar")
	if err == nil {
		t.Errorf("Get did not return an error with invalid id: foobar")
	}
	if !errors.Is(err, models.ErrNoRecord) {
		t.Errorf("Get returned an unexpected error: %s", err.Error())
	}

	if result != nil {
		t.Errorf("Get did not return a nil result: %v", result)
	}

}

func TestGetAllEmpty(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)

	results := db.GetAll()
	if len(results) != 0 {
		t.Errorf("expected 0 records, got %d : %v", len(results), results)
	}
}

func TestGetAll(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)

	for i := 0; i < 100; i++ {
		_, err := db.Add(models.Telemetry{Id: uuid.New().String()})
		if err != nil {
			t.Errorf("error adding record to db: %s", err.Error())
		}
	}

	results := db.GetAll()
	if len(results) != 100 {
		t.Errorf("expected 100 records, got %d : %v", len(results), results)
	}

}

func TestAliveCheck(t *testing.T) {
	// the in-memory db always returns a postive alive check
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)

	data, err := db.Alive()
	if err != nil {
		t.Errorf("alive health check returned an error: %s", err.Error())
	}

	if _, ok := data["health"]; !ok {
		t.Errorf("alive health check missing health key")
	}
}

func TestReadyCheck(t *testing.T) {
	// the in-memory db always returns a postive ready check
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	db := New(&logger)

	data, err := db.Ready()
	if err != nil {
		t.Errorf("ready health check returned an error: %s", err.Error())
	}

	if _, ok := data["ready"]; !ok {
		t.Errorf("ready health check missing ready key")
	}
}
