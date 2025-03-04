// Package conf provides and implements the interface to manage the Configurations.
package conf

import (
	"context"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/spf13/cast"
)

// Conf is a registry interface
type Conf interface {
	// WithReaders stores the given readers to load the data in the Load function
	WithReaders(readers ...Reader) Conf
	// WithTransformers stores the given transformers to change the output of the Get function.
	// All transformers will be applied in the given order.
	WithTransformers(transformers ...Transform) Conf

	// Reset creates an empty storage and clears the old one
	// It does not clear the default values
	// Can be used in the Unit Tests
	Reset() Conf
	// Load calls the `Read` function of all readers  in provided order
	Load(ctx context.Context) error
	// Keys returns the list of the stored keys
	Keys() []string
	// SetDefault sets a default value for a key
	SetDefault(key string, value interface{}) Conf
	// Set overrides the current value of a given key.
	Set(key string, value interface{}) Conf
	// Get returns a value for a given key if it is set or default value
	// Returns `nil` if key not found

	Get(key string) interface{}
	// GetString casts a value for a given key to String
	GetString(key string) string
	// GetInt casts a value for a given key to Int
	GetInt(key string) int
	// GetInt8 casts a value for a given key to Int8
	GetInt8(key string) int8
	// GetInt16 casts a value for a given key to Int16
	GetInt16(key string) int16
	// GetInt32 casts a value for a given key to Int32
	GetInt32(key string) int32
	// GetInt64 casts a value for a given key to Int64
	GetInt64(key string) int64
	// GetBool casts a value for a given key to Bool
	GetBool(key string) bool
	// GetFloat32 casts a value for a given key to Float32
	GetFloat32(key string) float32
	// GetFloat64 casts a value for a given key to Float64
	GetFloat64(key string) float64
	// GetTime casts a value for a given key to `time.Time`
	GetTime(key string) time.Time
	// GetDuration casts a value for a given key to `time.Duration`
	GetDuration(key string) time.Duration
}

type conf struct {
	storage  *sync.Map
	defaults *sync.Map

	readers      []Reader
	transformers []Transform
}

// New crates an instance of Conf interface
func New() Conf {
	c := &conf{
		storage:  &sync.Map{},
		defaults: &sync.Map{},
	}
	return c
}

var globalConf = New()

// GlobalConf returns the global Conf object
func GlobalConf() Conf {
	return globalConf
}

// WithReaders overrides the readers with the given ones
// The alias to work with an instance of the global configuration manager.
func WithReaders(readers ...Reader) Conf {
	return globalConf.WithReaders(readers...)
}

func (c *conf) WithReaders(readers ...Reader) Conf {
	c.readers = readers
	return c
}

// WithTransformers stores the given transformers to change the output of the Get function
// All transformers will be applied in the given order.
// The alias to work with an instance of the global configuration manager.
func WithTransformers(transformers ...Transform) Conf {
	return globalConf.WithTransformers(transformers...)
}

func (c *conf) WithTransformers(transformers ...Transform) Conf {
	c.transformers = transformers
	return c
}

// Reset creates an empty storage and clears the old one
// It does not clear the default values
// Can be used in the Unit Tests
// The alias to work with an instance of the global configuration manager.
func Reset() Conf {
	return globalConf.Reset()
}

func (c *conf) Reset() Conf {
	old := atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&c.storage)), //nolint:gosec
		unsafe.Pointer(&sync.Map{}),                   //nolint:gosec
	)
	s := (*sync.Map)(old)

	var keys []interface{}
	s.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	for _, key := range keys {
		s.Delete(key)
	}

	return c
}

// Load calls the `Read` function of all readers  in provided order
// The alias to work with an instance of the global configuration manager.
func Load(ctx context.Context) error {
	return globalConf.Load(ctx)
}

func (c *conf) Load(ctx context.Context) error {
	c.Reset()

	for _, reader := range c.readers {
		data, err := reader.Read(ctx)
		if err != nil {
			return err
		}

		c.scan(data, reader.Prefix())
	}

	return nil
}

func (c *conf) scan(data interface{}, key string) {
	if key != "" {
		c.storage.Store(key, data)
		key += "."
	}

	v := reflect.ValueOf(data)

	switch v.Kind() { //nolint:exhaustive
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			c.scan(iter.Value().Interface(), key+iter.Key().String())
		}
	case reflect.Array, reflect.Slice:
		for i := range v.Len() {
			c.scan(v.Index(i).Interface(), key+strconv.Itoa(i))
		}
	default:
	}
}

// Keys returns the list of the stored keys
// The alias to work with an instance of the global configuration manager.
func Keys() []string {
	return globalConf.Keys()
}

func (c *conf) Keys() []string {
	var keys []string

	c.storage.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	c.defaults.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	return keys
}

// SetDefault sets a default value for a key
// The alias to work with an instance of the global configuration manager.
func SetDefault(key string, value interface{}) Conf {
	return globalConf.SetDefault(key, value)
}

func (c *conf) SetDefault(key string, value interface{}) Conf {
	c.defaults.Store(key, value)
	return c
}

// Set overrides the current value of a given key.
// The alias to work with an instance of the global configuration manager.
func Set(key string, value interface{}) Conf {
	return globalConf.Set(key, value)
}

func (c *conf) Set(key string, value interface{}) Conf {
	c.storage.Store(key, value)
	return c
}

// Get returns a value for a given key if it is set or default value
// Returns `nil` if key not found
// The alias to work with an instance of the global configuration manager.
func Get(key string) interface{} {
	return globalConf.Get(key)
}

func (c *conf) Get(key string) interface{} {
	value, ok := c.storage.Load(key)
	if !ok {
		value, _ = c.defaults.Load(key)
	}

	for _, tr := range c.transformers {
		value = tr(key, value, c)
	}

	return value
}

// GetString casts a value for a given key to String
// The alias to work with an instance of the global configuration manager.
func GetString(key string) string {
	return globalConf.GetString(key)
}

func (c *conf) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// GetInt casts a value for a given key to Int
// The alias to work with an instance of the global configuration manager.
func GetInt(key string) int {
	return globalConf.GetInt(key)
}

func (c *conf) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// GetInt8 casts a value for a given key to Int8
// The alias to work with an instance of the global configuration manager.
func GetInt8(key string) int8 {
	return globalConf.GetInt8(key)
}

func (c *conf) GetInt8(key string) int8 {
	return cast.ToInt8(c.Get(key))
}

// GetInt16 casts a value for a given key to Int16
// The alias to work with an instance of the global configuration manager.
func GetInt16(key string) int16 {
	return globalConf.GetInt16(key)
}

func (c *conf) GetInt16(key string) int16 {
	return cast.ToInt16(c.Get(key))
}

// GetInt32 casts a value for a given key to Int32
// The alias to work with an instance of the global configuration manager.
func GetInt32(key string) int32 {
	return globalConf.GetInt32(key)
}

func (c *conf) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

// GetInt64 casts a value for a given key to Int64
// The alias to work with an instance of the global configuration manager.
func GetInt64(key string) int64 {
	return globalConf.GetInt64(key)
}

func (c *conf) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

// BoolValues is a global extendable list of the string values that should be converted as true or false
//
// Example:
//
//	conf.GetSet("flag", "sí")
//	conf.GetBool("flag") == false
//	conf.TrueValues["sí"] = true
//	conf.GetBool("flag") == true
var BoolValues = map[string]bool{
	"yes": true,
	"Yes": true,
	"y":   true,
	"Y":   true,
	"On":  true,
	"on":  true,
}

// GetBool casts a value for a given key to Bool
// The alias to work with an instance of the global configuration manager.
func GetBool(key string) bool {
	return globalConf.GetBool(key)
}

func (c *conf) GetBool(key string) bool {
	value := c.Get(key)

	if v, err := cast.ToBoolE(value); err == nil {
		return v
	}

	if v, ok := value.(string); ok {
		return BoolValues[v]
	}

	return false
}

// GetFloat32 casts a value for a given key to Float32
// The alias to work with an instance of the global configuration manager.
func GetFloat32(key string) float32 {
	return globalConf.GetFloat32(key)
}

func (c *conf) GetFloat32(key string) float32 {
	return cast.ToFloat32(c.Get(key))
}

// GetFloat64 casts a value for a given key to Float64
// The alias to work with an instance of the global configuration manager.
func GetFloat64(key string) float64 {
	return globalConf.GetFloat64(key)
}

func (c *conf) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// GetTime casts a value for a given key to `time.Time`
// The alias to work with an instance of the global configuration manager.
func GetTime(key string) time.Time {
	return globalConf.GetTime(key)
}

func (c *conf) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

// GetDuration casts a value for a given key to `time.Duration`
// The alias to work with an instance of the global configuration manager.
func GetDuration(key string) time.Duration {
	return globalConf.GetDuration(key)
}

func (c *conf) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}
