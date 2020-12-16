package inmem

import (
	"fmt"
	"time"

	"github.com/nedscode/memdb"
	"github.com/rs/zerolog"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

type InMemoryDB struct {
	db  *memdb.Store
	log *zerolog.Logger
}

func New(logger *zerolog.Logger) *InMemoryDB {
	mdb := memdb.NewStore().PrimaryKey("id").Unique()
	return &InMemoryDB{
		db:  mdb,
		log: logger,
	}
}

// Expire will expire all objects that have exceeded their TTL
// If successful, Expire will return the number of objects expired.
func (mem *InMemoryDB) Expire() int {
	var count int
	objects := mem.GetAll()
	for _, obj := range objects {
		if obj.IsExpired() {
			mem.log.Debug().Str("obj", obj.Id).Msg("object telemetry is stale")
			mem.db.Delete(obj)
			count++
		}
	}
	mem.log.Info().Int("objects", count).Msg("stale objects expired")
	return count
}

// Add a new telemetry struct to the in memory database and return its id
// as a string.  If the object can't be added, return an error.
func (mem *InMemoryDB) Add(t models.Telemetry) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		models.TransactionDuration.WithLabelValues("inmemdb", "Add").Observe(duration.Seconds())
	}()

	old, err := mem.db.Put(&t)
	if err != nil {
		models.TransactionErrors.WithLabelValues("inmemdb", "Add").Inc()
		return "", err
	}

	if old == nil {
		// must be a new record
		models.RecordCount.WithLabelValues("inmemdb").Inc()
	}

	return t.Id, nil
}

// Get will return the telemetry of the object with the passed id.
// If the object is not found then a NotFound error is returned.
func (mem *InMemoryDB) Get(id string) (*models.Telemetry, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		models.TransactionDuration.WithLabelValues("inmemdb", "Get").Observe(duration.Seconds())
	}()

	found := mem.db.InPrimaryKey().One(id)
	if t, ok := found.(*models.Telemetry); ok {
		return t, nil
	}
	models.TransactionErrors.WithLabelValues("inmemdb", "Get").Inc()
	return nil, models.ErrNoRecord
}

// GetAll will return all known telemtry objects
func (mem *InMemoryDB) GetAll() []models.Telemetry {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		models.TransactionDuration.WithLabelValues("inmemdb", "GetAll").Observe(duration.Seconds())
	}()
	var results []models.Telemetry
	mem.db.Ascend(func(indexer interface{}) bool {
		t := indexer.(*models.Telemetry)
		results = append(results, *t)
		return true
	})

	return results
}

// Alive returns the health status of the database
// If the database is in a state the is nonrecoverable it will
// return an error
func (mem *InMemoryDB) Alive() (map[string]string, error) {
	return map[string]string{
		"health": "alive",
		"items":  fmt.Sprintf("%d", mem.db.Len()),
	}, nil
}

// Reader will return the health status of the database
// An error will be returned if the database should not
// currently accepted new connections
func (mem *InMemoryDB) Ready() (map[string]string, error) {
	return map[string]string{
		"health":  "alive",
		"ready":   "true",
		"items":   fmt.Sprintf("%d", mem.db.Len()),
		"message": fmt.Sprintf("up; %d active objects", mem.db.Len()),
	}, nil
}
