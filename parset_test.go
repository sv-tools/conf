package conf_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sv-tools/conf"
)

func testParseFunc(_ context.Context, r io.Reader) (any, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	for _, line := range strings.Split(string(data), ";") {
		values := strings.Split(line, ":")
		res[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
	}

	return res, nil
}

func testParseFuncError(_ context.Context, _ io.Reader) (any, error) {
	return nil, errFake
}

func TestStreamParser(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`foo:1;bar:2`))
	c := conf.New().WithReaders(conf.NewStreamParser(reader).WithParser(testParseFunc).WithPrefix("pr"))
	require.NoError(t, c.Load(t.Context()))
	require.Equal(t, 1, c.GetInt("pr.foo"))
	require.Equal(t, 2, c.GetInt("pr.bar"))
}

func TestStreamParser_ErrNoParser(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`foo:1;bar:2`))
	c := conf.New().WithReaders(conf.NewStreamParser(reader))
	require.ErrorIs(t, c.Load(t.Context()), conf.ErrNoParser)
}

func TestStreamParser_ErrNoStream(t *testing.T) {
	t.Parallel()

	c := conf.New().WithReaders(conf.NewStreamParser(nil).WithParser(testParseFunc))
	require.ErrorIs(t, c.Load(t.Context()), conf.ErrNoStream)
}

func TestStreamParser_errFake(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte(`foo:1;bar:2`))
	c := conf.New().WithReaders(conf.NewStreamParser(reader).WithParser(testParseFuncError))
	require.ErrorIs(t, c.Load(t.Context()), errFake)
}

func TestFileParser(t *testing.T) {
	t.Parallel()

	parser, err := conf.NewFileParser(`testdata/data.txt`)
	require.NoError(t, err)

	c := conf.New().WithReaders(parser.WithParser(testParseFunc).WithPrefix("pr"))
	require.NoError(t, c.Load(t.Context()))
	require.Equal(t, 1, c.GetInt("pr.foo"))
	require.Equal(t, 2, c.GetInt("pr.bar"))
}

func TestFileParser_FileNotFound(t *testing.T) {
	t.Parallel()

	parser, err := conf.NewFileParser(`testdata/fake.txt`)
	require.EqualError(t, err, "open testdata/fake.txt: no such file or directory")
	require.Nil(t, parser)
}
