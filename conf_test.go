package conf_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sv-tools/conf"
)

type testReader struct {
	err    error
	data   interface{}
	prefix string
}

func (t *testReader) Read(_ context.Context) (interface{}, error) {
	return t.data, t.err
}

func (t *testReader) Prefix() string {
	return t.prefix
}

func newReader(tb testing.TB, prefix string, data interface{}, err error) conf.Reader {
	tb.Helper()

	return &testReader{prefix: prefix, data: data, err: err}
}

var errFake = errors.New("fake error")

func TestConf(t *testing.T) {
	t.Parallel()

	data1 := map[string]interface{}{
		"foo": "bar",
		"baz": 42,
		"xyz": []int{1, 2, 3},
		"a": map[string]interface{}{
			"b": []map[string]interface{}{
				{
					"c": 1,
				},
				{
					"d": 2,
				},
			},
		},
		"x.y.z": "xyz",
	}
	data2 := "fake data"
	data3 := 34
	data4 := []interface{}{
		1,
		"2",
	}

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		require.EqualError(t,
			conf.New().WithReaders(newReader(t, "", nil, errFake)).Load(context.Background()),
			"fake error",
		)
	})

	t.Run("full", func(t *testing.T) {
		t.Parallel()

		c := conf.New().WithReaders(
			newReader(t, "", data1, nil),
			newReader(t, "data2", data2, nil),
			newReader(t, "data3", data3, nil),
			newReader(t, "data4", data4, nil),
		)
		require.NoError(t, c.Load(context.Background()))

		require.Nil(t, c.Get("no key"))
		require.Equal(t, data4, c.Get("data4"))
		require.Equal(t, 1, c.Get("data4.0"))
		require.Equal(t, "2", c.Get("data4.1"))

		require.Equal(t, data3, c.Get("data3"))
		require.Nil(t, c.Get("data3.no"))

		require.Equal(t, data2, c.Get("data2"))
		require.Nil(t, c.Get("data2.no"))

		require.Equal(t, "bar", c.Get("foo"))
		require.Equal(t, 42, c.Get("baz"))
		require.Equal(t, data1["xyz"], c.Get("xyz"))
		require.Equal(t, 2, c.Get("xyz.1"))
		require.Equal(t, 1, c.Get("a.b.0.c"))
		require.Equal(t, 2, c.Get("a.b.1.d"))
		require.Nil(t, c.Get("a.b.e.no"))
		require.Equal(t, "xyz", c.Get("x.y.z"))
		require.Nil(t, c.Get("x.y"))
	})

	t.Run("list only", func(t *testing.T) {
		t.Parallel()

		c := conf.New().WithReaders(
			newReader(t, "", data4, nil),
		)
		require.NoError(t, c.Load(context.Background()))
		for i, expectedValue := range data4 {
			require.Equal(t, expectedValue, c.Get(strconv.Itoa(i)))
		}
	})

	t.Run("single value without prefix", func(t *testing.T) {
		t.Parallel()

		c := conf.New().WithReaders(
			newReader(t, "", data3, nil),
		)
		require.NoError(t, c.Load(context.Background()))
		require.Empty(t, c.Keys())
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		c := conf.New()
		c.SetDefault("foo", 42)
		require.Equal(t, 42, c.Get("foo"))
		c.Set("foo", 101)
		require.Equal(t, 101, c.Get("foo"))

		c.Reset()
		require.Equal(t, 42, c.Get("foo"))
	})

	t.Run("keys", func(t *testing.T) {
		t.Parallel()

		c := conf.New()
		require.Empty(t, c.Keys())
		c.Set("bar", 101)
		require.Equal(t, []string{"bar"}, c.Keys())

		c.SetDefault("foo", 42)
		require.ElementsMatch(t, []string{"foo", "bar"}, c.Keys())

		c.Reset()
		require.Equal(t, []string{"foo"}, c.Keys())
	})
}

func TestConf_GetBool(t *testing.T) {
	t.Parallel()

	c := conf.New()
	data := map[interface{}]bool{
		"0":     false,
		"1":     true,
		0:       false,
		2:       true,
		true:    true,
		false:   false,
		"":      false,
		"True":  true,
		"False": false,
		"yes":   true,
		"no":    false,
		"foo":   false,
		"sí":    false,
		"34":    false,
		34:      true,
		1.25:    true,
		1.0:     true,
		nil:     false,
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetBool("flag"), "%T: %v", rawValue, rawValue)
	}

	conf.BoolValues["sí"] = true
	c.Set("flag", "sí")
	require.True(t, c.GetBool("flag"))
}

func TestConf_GetString(t *testing.T) {
	t.Parallel()

	c := conf.New()
	data := map[interface{}]string{
		"0":   "0",
		"1":   "1",
		0:     "0",
		-1:    "-1",
		true:  "true",
		false: "false",
		"":    "",
		1.25:  "1.25",
		nil:   "",
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetString("flag"), "%T: %v", rawValue, rawValue)
	}
}

func TestConf_GetInt(t *testing.T) {
	t.Parallel()

	c := conf.New()
	data := map[interface{}]int{
		"0":   0,
		"-1":  -1,
		0:     0,
		-1:    -1,
		true:  1,
		false: 0,
		"":    0,
		1.25:  1,
		1.99:  1,
		nil:   0,
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetInt("flag"), "%T: %v", rawValue, rawValue)
		require.Equal(t, int8(expectedValue), c.GetInt8("flag"), "%T: %v", rawValue, rawValue)
		require.Equal(t, int16(expectedValue), c.GetInt16("flag"), "%T: %v", rawValue, rawValue)
		require.Equal(t, int32(expectedValue), c.GetInt32("flag"), "%T: %v", rawValue, rawValue)
		require.Equal(t, int64(expectedValue), c.GetInt64("flag"), "%T: %v", rawValue, rawValue)
	}
}

func TestConf_GetFloat(t *testing.T) {
	t.Parallel()

	c := conf.New()
	data := map[interface{}]float64{
		"0.1":  0.1,
		"-1.4": -1.4,
		0:      0,
		-1:     -1,
		true:   1,
		false:  0,
		"":     0,
		1.25:   1.25,
		1.99:   1.99,
		nil:    0,
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetFloat64("flag"), "%T: %v", rawValue, rawValue)
		require.Equal(t, float32(expectedValue), c.GetFloat32("flag"), "%T: %v", rawValue, rawValue)
	}
}

func TestConf_GetTime(t *testing.T) {
	t.Parallel()

	t1, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	require.NoError(t, err)

	c := conf.New()
	data := map[interface{}]time.Time{
		"0.1":                  {},
		"-1.4":                 {},
		0:                      time.Unix(0, 0),
		-1:                     time.Unix(-1, 0),
		25:                     time.Unix(25, 0),
		"25":                   {},
		true:                   {},
		false:                  {},
		"":                     {},
		1.25:                   {},
		1.99:                   {},
		nil:                    {},
		"2006-01-02T15:04:05Z": t1,
		"2h1m25s":              {},
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetTime("flag"), "%T: %v", rawValue, rawValue)
	}
}

func TestConf_GetDuration(t *testing.T) {
	t.Parallel()

	c := conf.New()
	data := map[interface{}]time.Duration{
		"0.1":                  0,
		"-1.4":                 -1,
		0:                      0,
		-1:                     -1,
		25:                     25,
		"25":                   25,
		true:                   0,
		false:                  0,
		"":                     0,
		1.25:                   1,
		1.99:                   1,
		nil:                    0,
		"2006-01-02T15:04:05Z": 0,
		"1m":                   60000000000,
		"2h1m25s":              7285000000000,
	}
	for rawValue, expectedValue := range data {
		c.Set("flag", rawValue)

		require.Equal(t, expectedValue, c.GetDuration("flag"), "%T: %v", rawValue, rawValue)
	}
}

func TestGlobalConf(t *testing.T) {
	t.Cleanup(func() {
		conf.Reset()
	})

	require.Nil(t, conf.Get("foo"))
	conf.WithReaders(newReader(t, "", map[string]interface{}{
		"foo": 42,
	}, nil))
	require.NoError(t, conf.Load(context.Background()))
	require.Equal(t, 42, conf.Get("foo"))

	require.Nil(t, conf.Get("bar"))
	conf.Set("bar", 42)
	require.Equal(t, 42, conf.Get("bar"))

	require.Nil(t, conf.Get("xyz"))
	conf.GlobalConf().Set("xyz", 101)
	require.Equal(t, 101, conf.Get("xyz"))

	require.Nil(t, conf.Get("default"))
	conf.SetDefault("default", 101)
	require.Equal(t, 101, conf.Get("default"))

	require.True(t, conf.GetBool("foo"))
	require.False(t, conf.GetBool("no"))
	require.Equal(t, 42, conf.GetInt("foo"))
	require.Equal(t, int8(42), conf.GetInt8("foo"))
	require.Equal(t, int16(42), conf.GetInt16("foo"))
	require.Equal(t, int32(42), conf.GetInt32("foo"))
	require.Equal(t, int64(42), conf.GetInt64("foo"))
	require.Equal(t, float32(42), conf.GetFloat32("foo"))
	require.Equal(t, float64(42), conf.GetFloat64("foo"))
	require.Equal(t, "42", conf.GetString("foo"))
	require.Equal(t, time.Unix(42, 0), conf.GetTime("foo"))
	require.Equal(t, time.Duration(42), conf.GetDuration("foo"))

	require.ElementsMatch(t, []string{"foo", "bar", "xyz", "default"}, conf.Keys())
}

func testTransform(_ string, value interface{}, _ conf.Conf) interface{} {
	if v, ok := value.(string); ok && v == "value-to-be-transformed" {
		return 101
	}

	return value
}

func TestConf_WithTransformer(t *testing.T) {
	t.Parallel()

	c := conf.New().WithTransformers(testTransform)
	c.Set("foo", "33")
	c.Set("bar", "value-to-be-transformed")

	require.Equal(t, "33", c.Get("foo"))
	require.Equal(t, 101, c.Get("bar"))
}

func TestGlobalConf_WithTransformer(t *testing.T) {
	t.Cleanup(func() {
		conf.Reset()
	})

	c := conf.WithTransformers(testTransform)
	c.Set("foo", "33")
	c.Set("bar", "value-to-be-transformed")

	require.Equal(t, "33", c.Get("foo"))
	require.Equal(t, 101, c.Get("bar"))
}
