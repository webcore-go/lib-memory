package redis

import (
	"fmt"
	"reflect"
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

	// Handle pointer types by dereferencing
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("Cannot set nil pointer value for key %s", key)
		}
		value = rv.Elem().Interface()
	}

	switch v := value.(type) {
	case string:
		val = v
	case int:
		val = strconv.FormatInt(int64(v), 10)
	case int8:
		val = strconv.FormatInt(int64(v), 10)
	case int16:
		val = strconv.FormatInt(int64(v), 10)
	case int32:
		val = strconv.FormatInt(int64(v), 10)
	case int64:
		val = strconv.FormatInt(v, 10)
	case uint:
		val = strconv.FormatUint(uint64(v), 10)
	case uint8:
		val = strconv.FormatUint(uint64(v), 10)
	case uint16:
		val = strconv.FormatUint(uint64(v), 10)
	case uint32:
		val = strconv.FormatUint(uint64(v), 10)
	case uint64:
		val = strconv.FormatUint(v, 10)
	case bool:
		val = "0"
		if v {
			val = "1"
		}
	case float32:
		val = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		val = strconv.FormatFloat(v, 'f', -1, 64)
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

func (r *MemoryCache) Get(key string, outvalue any) bool {
	val, ok := r.Cache.GetIfPresent(key)
	if !ok {
		return false
	}

	// Use reflection to properly set the outvalue
	rv := reflect.ValueOf(outvalue)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}

	elem := rv.Elem()

	switch elem.Kind() {
	case reflect.String:
		elem.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			elem.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			elem.SetUint(i)
		}
	case reflect.Bool:
		elem.SetBool(val == "1")
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err == nil {
			elem.SetFloat(f)
		}
	default:
		// For complex types, use JSON unmarshal
		helper.JSONUnmarshal([]byte(val), outvalue)
	}

	return ok
}
