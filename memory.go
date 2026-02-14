package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/webcore-go/webcore/app/helper"
	"github.com/webcore-go/webcore/infra/config"

	"github.com/maypok86/otter/v2"
	"github.com/maypok86/otter/v2/stats"
)

// MemoryCache represents shared MemoryCache connection
type MemoryCache struct {
	Cache *otter.Cache[string, string]
}

// NewMemoryCache creates a new MemoryCache connection
func NewMemoryCache(config config.MemoryConfig) (*MemoryCache, error) {
	// Create statistics counter to track cache operations
	counter := stats.NewCounter()

	limit := 10_000
	if config.Limit > 0 {
		limit = config.Limit
	}

	cache := otter.Must(&otter.Options[string, string]{
		MaximumSize:       limit,
		ExpiryCalculator:  otter.ExpiryAccessing[string, string](config.ExpiresIn),      // Reset timer on reads/writes
		RefreshCalculator: otter.RefreshWriting[string, string](500 * time.Millisecond), // Refresh after writes
		StatsRecorder:     counter,                                                      // Attach stats collector
	})

	return &MemoryCache{Cache: cache}, nil
}

func (r *MemoryCache) Install(args ...any) error {
	// Tidak melakukan apa-apa
	return nil
}

func (r *MemoryCache) Connect() error {
	// Tidak melakukan apa-apa
	return nil
}

// Close closes the MemoryCache connection
func (r *MemoryCache) Disconnect() error {
	// Tidak melakukan apa-apa
	return nil
}

func (r *MemoryCache) Uninstall() error {
	// Tidak melakukan apa-apa
	return nil
}

func (r *MemoryCache) Set(key string, value any, ttl time.Duration) error {
	var val string
	var ok bool

	switch v := value.(type) {
	case string:
		val = v
	case int:
	case int16:
	case int64:
		val = strconv.FormatInt(int64(v), 64)
	case bool:
		val = "0"
		if v {
			val = "1"
		}
	case float32:
	case float64:
		val = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		m, err := helper.JSONMarshal(value)
		if err != nil {
			return err
		}

		val = string(m)
	}

	_, ok = r.Cache.Set(key, val)
	if !ok {
		return fmt.Errorf("Gagal simpan di MemoryCache %s", key)
	}
	return nil
}

func (r *MemoryCache) Get(key string) (any, bool) {
	val, ok := r.Cache.GetIfPresent(key)
	var val2 any
	val2 = val
	return val2, ok
}
